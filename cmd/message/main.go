package main

import (
	"github.com/Vilinvil/task_messaggio/internal/message/config"
	"github.com/Vilinvil/task_messaggio/internal/message/server"
)

func main() {
	configServer, err := config.New()
	if err != nil {
		panic(err)
	}

	srv := new(server.Server)

	if err := srv.Run(configServer); err != nil {
		panic(err)
	}
}
