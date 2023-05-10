package config

import (
	"fmt"
	"os"
)

const (
	LarkEnvKeyAppId     = "lark_app_id"
	LarkEnvKeyAppSecret = "lark_app_secret"
)

const (
	PaginationStorageBitableAppToken = "DLSbbQIcEa0KyIsetHWcg3PDnNh"
	PaginationStorageBitableTableId  = "tblLJY5YSoEkV3G3"
	ChatStorageBitableAppToken       = "QCBrbzgx4aKRAis9eewcV731n7d"
	ChatStorageBitableTableId        = "tblIk692K5LXte8x"
)

type LarkConfig struct {
	AppId     string
	AppSecret string
}

// NewLarkConfig lookup lark configs from the environment variables,
// otherwise panic.
func NewLarkConfig() LarkConfig {
	return LarkConfig{
		AppId:     lookupEnvVariableFor(LarkEnvKeyAppId),
		AppSecret: lookupEnvVariableFor(LarkEnvKeyAppSecret),
	}
}

func lookupEnvVariableFor(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("environment variable %v is not defined.", key))
	}
	if len(value) == 0 {
		panic(fmt.Sprintf("environment variable %v should not be empty.", key))
	}

	return value
}
