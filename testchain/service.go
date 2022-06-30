package testchain

import "time"

type Request struct {
	message string
	t       time.Time
	err     error
}

type Response struct {
	status int
	err    error
}

type Trigger interface {
	TurnOn()
	TurnOff()
}

type Service struct {
	storage   map[Request]Response
	reachable bool
	Trigger
}

func MakeService() *Service {
	r := true
	s := make(map[Request]Response)
	srv := &Service{storage: s, reachable: r}
	return srv
}

func (s *Service) TurnOn() {
	s.reachable = true
}

func (s *Service) TurnOff() {
	s.reachable = false
}
