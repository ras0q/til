package mock

import (
	"log"

	"github.com/ras0q/go-playground-test-overlay/caller"
)

type CallerMock struct{}

// interface guard
var _ caller.Caller = (*CallerMock)(nil)

func New() *CallerMock {
	return &CallerMock{}
}

var msg = "called"

func (s *CallerMock) Call() string {
	log.Println(msg)

	return msg
}

// extra method
func (s *CallerMock) SetMsg(m string) {
	msg = m
}
