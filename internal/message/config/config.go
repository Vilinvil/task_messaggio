package config

import (
	"fmt"
	"strconv"
	"time"

	"github.com/Vilinvil/task_messaggio/pkg/config"
	"github.com/Vilinvil/task_messaggio/pkg/myerrors"
)

const (
	EnvPort = "PORT"

	StandardPort               = "8090"
	StandardOutputLogPath      = "stdout /var/log/message/logs.json"
	StandardErrorOutputLogPath = "stderr /var/log/message/err_logs.json"
)

type Config struct {
	Port               string
	OutputLogPath      string
	ErrorOutputLogPath string
	BasicTimeout       time.Duration
	URLDataBase        string
	APIName            string
}

func New() (*Config, error) {
	timeoutInSecond, err := strconv.Atoi(config.GetEnvStr(config.EnvBasicTimeout, config.StandardBasicTimeout))
	if err != nil {
		return nil, fmt.Errorf(myerrors.ErrTemplate, err)
	}

	basicTimeout := time.Duration(timeoutInSecond) * time.Second

	return &Config{
		Port:               config.GetEnvStr(EnvPort, StandardPort),
		OutputLogPath:      config.GetEnvStr(config.EnvOutputLogPath, StandardOutputLogPath),
		ErrorOutputLogPath: config.GetEnvStr(config.EnvErrorOutputLogPath, StandardErrorOutputLogPath),
		BasicTimeout:       basicTimeout,
		URLDataBase:        config.GetEnvStr(config.EnvURLDataBase, config.StandardURLDataBase),
		APIName:            config.GetEnvStr(config.EnvAPIName, config.StandardAPIName),
	}, nil
}
