package config

import (
	"os"
)

const (
	EnvOutputLogPath      = "OUTPUT_LOG_PATH"
	EnvErrorOutputLogPath = "ERROR_OUTPUT_LOG_PATH"
	EnvBasicTimeout
	EnvURLDataBase = "URL_DATABASE"
	EnvAPIName     = "API_NAME"

	StandardBasicTimeout = "10"
	StandardURLDataBase  = "postgres://username:wrongpassword@localhost:5432/dbname?sslmode=disable"
	StandardAPIName      = "/api/v1"
)

func GetEnvStr(envName string, standardValue string) string {
	result, ok := os.LookupEnv(envName)
	if !ok {
		return standardValue
	}

	return result
}
