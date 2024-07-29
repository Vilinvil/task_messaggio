package config

import (
	"os"
)

const (
	EnvOutputLogPath      = "OUTPUT_LOG_PATH"
	EnvErrorOutputLogPath = "ERROR_OUTPUT_LOG_PATH"
	EnvBasicTimeout       = "BASIC_TIMEOUT"
	EnvURLDataBase        = "URL_DATABASE"
	EnvAPIName            = "API_NAME"
	EnvProductionMode     = "PRODUCTION_MODE"
	EnvBrokerAddr         = "BROKER_ADDR"

	StandardBasicTimeout   = "10"
	StandardURLDataBase    = "postgres://username:wrongpassword@localhost:5432/dbname?sslmode=disable"
	StandardAPIName        = "/api/v1"
	StandardProductionMode = false
	StandardBrokerAddr     = "localhost:9092"
)

func GetEnvStr(envName string, standardValue string) string {
	result, ok := os.LookupEnv(envName)
	if !ok {
		return standardValue
	}

	return result
}

func GetEnvBool(envName string, standardValue bool) bool {
	result, ok := os.LookupEnv(envName)
	if !ok {
		return standardValue
	}

	return result == "true"
}
