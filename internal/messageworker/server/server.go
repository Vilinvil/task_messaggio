package server

import (
	"github.com/Vilinvil/task_messaggio/internal/messageworker/messageworker/usecases"
	"github.com/Vilinvil/task_messaggio/pkg/mylogger"
)

var _ MessageWorker = (*usecases.SleepingMessageWorker)(nil)

type MessageWorker interface {
	JobMessages() error
}

func Run(logger *mylogger.MyLogger) *ServerWorker {

}
