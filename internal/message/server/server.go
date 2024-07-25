package server

import (
	"context"
	"fmt"
	"github.com/Vilinvil/task_messaggio/internal/message/server/delivery"
	"net/http"
	"strings"

	"github.com/Vilinvil/task_messaggio/internal/message/config"
	"github.com/Vilinvil/task_messaggio/pkg/myerrors"
	"github.com/Vilinvil/task_messaggio/pkg/mylogger"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(config *config.Config) (err error) {
	baseCtx := context.Background()

	logger, err := mylogger.New(strings.Split(config.OutputLogPath, " "),
		strings.Split(config.ErrorOutputLogPath, " "))
	if err != nil {
		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	defer logger.Sync()

	mux, err := delivery.NewMux(baseCtx, logger)
	if err != nil {
		return fmt.Errorf(myerrors.ErrTemplate, err)
	}

	s.httpServer = &http.Server{ //nolint:exhaustruct
		Addr:           ":" + config.Port,
		Handler:        mux,
		MaxHeaderBytes: http.DefaultMaxHeaderBytes,
		ReadTimeout:    config.BasicTimeout,
		WriteTimeout:   config.BasicTimeout,
	}

	logger.Infof("Start server:%s", config.Port)

	return fmt.Errorf(myerrors.ErrTemplate, s.httpServer.ListenAndServe())
}

func (s *Server) Shutdown(ctx context.Context) error {
	return fmt.Errorf(myerrors.ErrTemplate, s.httpServer.Shutdown(ctx))
}
