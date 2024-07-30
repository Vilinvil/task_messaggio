package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/Vilinvil/task_messaggio/internal/message/config"
	"github.com/Vilinvil/task_messaggio/internal/message/server/delivery"
	"github.com/Vilinvil/task_messaggio/pkg/mylogger"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(ctx context.Context, config *config.Config) error {
	logger, err := mylogger.New(strings.Split(config.OutputLogPath, " "),
		strings.Split(config.ErrorOutputLogPath, " "), config.ProductionMode)
	if err != nil {
		return err
	}

	defer logger.Sync() //nolint:errcheck

	mux, err := delivery.NewMux(ctx, config.URLDataBase, config.BrokerAddr, logger)
	if err != nil {
		return err
	}

	s.httpServer = &http.Server{ //nolint:exhaustruct
		Addr:           ":" + config.Port,
		Handler:        mux,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
		ReadTimeout:    config.BasicTimeout,
		WriteTimeout:   config.BasicTimeout,
	}

	logger.Infof("Start server:%s", config.Port)

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
