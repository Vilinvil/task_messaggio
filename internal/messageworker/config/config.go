package config

import (
	"strconv"
	"time"

	"github.com/Vilinvil/task_messaggio/pkg/config"
)

const (
	EnvTimeOnTask = "TIME_ON_TASK"

	StandardOutputLogPath      = "stdout /var/log/messageworker/logs.json"
	StandardErrorOutputLogPath = "stderr /var/log/messageworker/err_logs.json"
	StandardTimeOnTask         = "10"
)

type Config struct {
	ProductionMode     bool
	TimeOnTask         time.Duration
	OutputLogPath      string
	ErrorOutputLogPath string
	URLDataBase        string
	BrokerAddr         string
}

func New() (*Config, error) {
	timeOnTaskInSecond, err := strconv.Atoi(config.GetEnvStr(EnvTimeOnTask, StandardTimeOnTask))
	if err != nil {
		return nil, err
	}

	timeOnTask := time.Duration(timeOnTaskInSecond) * time.Second

	return &Config{
		ProductionMode:     config.GetEnvBool(config.EnvProductionMode, config.StandardProductionMode),
		TimeOnTask:         timeOnTask,
		OutputLogPath:      config.GetEnvStr(config.EnvOutputLogPath, StandardOutputLogPath),
		ErrorOutputLogPath: config.GetEnvStr(config.EnvErrorOutputLogPath, StandardErrorOutputLogPath),
		URLDataBase:        config.GetEnvStr(config.EnvURLDataBase, config.StandardURLDataBase),
		BrokerAddr:         config.GetEnvStr(config.EnvBrokerAddr, config.StandardBrokerAddr),
	}, nil
}
