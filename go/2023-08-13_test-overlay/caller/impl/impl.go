package impl

import (
	"log"

	"github.com/ras0q/go-playground-test-overlay/caller"
)

type CallerImpl struct{}

func New() caller.Caller {
	return &CallerImpl{}
}

var msg = "called"

func (s *CallerImpl) Call() string {
	log.Println(msg)

	return msg
}
