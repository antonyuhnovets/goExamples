package testchain

import (
	"context"
	"errors"
	"fmt"
	"time"
)

type Caller interface {
	SendRequest(context.Context, string) error
	RecieveResponse(context.Context, Response) error
}

type Applier interface {
	RecieveRequest(context.Context) error
	SendResponse(context.Context, Response) error
}

type Upstream struct {
	srv *Service
	Caller
}

type Downstream struct {
	srv *Service
	Applier
}

type Chain struct {
	up   *Upstream
	down *Downstream
}

func MakeChain() *Chain {
	down := &Downstream{srv: MakeService()}
	up := &Upstream{srv: MakeService()}

	return &Chain{up, down}
}

func (c *Chain) TurnOffDown() {
	c.down.srv.reachable = false
}

func (c *Chain) TurnOnDown() {
	c.down.srv.reachable = true
}

func (c *Chain) SendRequest(ctx context.Context, message string) error {
	var err error
	req := Request{message, time.Now(), err}

	if !c.down.srv.reachable {
		err = errors.New("Service unreachable")
		c.up.srv.storage[req] = Response{500, err}

		return err
	}

	c.up.srv.storage[req] = Response{}
	cont := context.WithValue(ctx, "request", req)
	c.RecieveRequest(cont)
	fmt.Println("Request sending successfull")
	return err
}

func (c *Chain) RecieveRequest(ctx context.Context) error {
	var err error
	var res Response

	req, ok := ctx.Value("request").(Request)
	if !ok {
		err = errors.New("Request reciever context failed")
		return err
	}

	switch req.message {
	case "GET", "POST":
		res = Response{200, nil}
		err = nil
	default:
		err = errors.New("bad request")
		res = Response{400, err}
	}

	cont := context.WithValue(ctx, "res", res)
	c.SendResponse(cont, res)

	fmt.Println("Request recieving successfull")
	return err
}

func (c *Chain) SendResponse(ctx context.Context, res Response) error {
	req, ok := ctx.Value("request").(Request)
	if !ok {
		return errors.New("Response sender context failed")
	}
	c.down.srv.storage[req] = res
	c.RecieveResponse(ctx, res)

	fmt.Println("Response sending successfull")
	return nil
}

func (c *Chain) RecieveResponse(ctx context.Context, res Response) error {
	cont, cancel := context.WithCancel(ctx)
	defer cancel()

	req, ok := cont.Value("request").(Request)
	if !ok {
		return errors.New("Response reciever context failed")
	}
	c.up.srv.storage[req] = res
	fmt.Println("Response recieving successfull")
	return nil
}
