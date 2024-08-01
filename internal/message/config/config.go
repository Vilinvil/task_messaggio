package config

import (
	"strconv"
	"time"

	"github.com/Vilinvil/task_messaggio/pkg/config"
)

const (
	EnvPort = "PORT"

	StandardPort               = "8090"
	StandardOutputLogPath      = "stdout /var/log/message/logs.json"
	StandardErrorOutputLogPath = "stderr /var/log/message/err_logs.json"
)

type Config struct {
	ProductionMode     bool
	BasicTimeout       time.Duration
	Port               string
	OutputLogPath      string
	ErrorOutputLogPath string
	URLDataBase        string
	APIName            string
	BrokerAddr         string
}

func New() (*Config, error) {
	timeoutInSecond, err := strconv.Atoi(config.GetEnvStr(config.EnvBasicTimeout, config.StandardBasicTimeout))
	if err != nil {
		return nil, err
	}

	basicTimeout := time.Duration(timeoutInSecond) * time.Second

	return &Config{
		ProductionMode:     config.GetEnvBool(config.EnvProductionMode, config.StandardProductionMode),
		BasicTimeout:       basicTimeout,
		Port:               config.GetEnvStr(EnvPort, StandardPort),
		OutputLogPath:      config.GetEnvStr(config.EnvOutputLogPath, StandardOutputLogPath),
		ErrorOutputLogPath: config.GetEnvStr(config.EnvErrorOutputLogPath, StandardErrorOutputLogPath),
		URLDataBase:        config.GetEnvStr(config.EnvURLDataBase, config.StandardURLDataBase),
		APIName:            config.GetEnvStr(config.EnvAPIName, config.StandardAPIName),
		BrokerAddr:         config.GetEnvStr(config.EnvBrokerAddr, config.StandardBrokerAddr),
	}, nil
}
