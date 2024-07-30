package main

import (
	"context"

	"github.com/Vilinvil/task_messaggio/internal/messageworker/config"
	"github.com/Vilinvil/task_messaggio/internal/messageworker/server"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	serverWorker := server.Server{}

	ctx := context.Background()

	if err := serverWorker.Run(ctx, cfg); err != nil {
		panic(err)
	}
}
