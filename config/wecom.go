package config

import (
	"fmt"
	"strconv"
	"time"
)

const (
	WeComEnvKeyCorpID                = "wecom_corp_id"
	WeComEnvKeyChatSyncSecret        = "wecom_chat_sync_secret"
	WeComEnvKeyChatSyncRsaPrivateKey = "wecom_chat_sync_rsa_private_key"
	WeComEnvKeyAgentID               = "wecom_agent_id"
	WeComEnvKeyAgentSecret           = "wecom_agent_secret"
)

const (
	WeComChatRecordSDKTimeout = 30 * time.Second
)

type WeComConfig struct {
	CorpID                string
	ChatSyncSecret        string
	ChatSyncRsaPrivateKey string
	AgentID               int64
	AgentSecret           string
}

func NewWeComConfig() WeComConfig {
	return WeComConfig{
		CorpID:                lookupEnvVariableFor(WeComEnvKeyCorpID),
		ChatSyncSecret:        lookupEnvVariableFor(WeComEnvKeyChatSyncSecret),
		ChatSyncRsaPrivateKey: lookupEnvVariableFor(WeComEnvKeyChatSyncRsaPrivateKey),
		AgentID:               int64From(WeComEnvKeyAgentID, lookupEnvVariableFor(WeComEnvKeyAgentID)),
		AgentSecret:           lookupEnvVariableFor(WeComEnvKeyAgentSecret),
	}
}

func int64From(key string, value string) int64 {
	result, err := strconv.ParseInt(value, 10, 64)

	if err != nil {
		panic(fmt.Sprintf("environment variable '%v' is not a int64 number, value: %v", key, value))
	}

	return result
}