package main

import (
	"context"

	"github.com/Vilinvil/task_messaggio/internal/message/config"
	"github.com/Vilinvil/task_messaggio/internal/message/server"
)

// @title           Swagger message API
// @version         1.0
// @description     This api for message service
// @contact.name   Vladislav
// @contact.url    https://t.me/Vilin0
// @contact.email  ivn-15-07@mail.ru
// @host      goods-galaxy.ru:7443
// @BasePath  /api/v1
// @securityDefinitions.basic  BasicAuth
func main() {
	configServer, err := config.New()
	if err != nil {
		panic(err)
	}

	srv := new(server.Server)
	ctx := context.Background()

	if err := srv.Run(ctx, configServer); err != nil {
		panic(err)
	}
}
