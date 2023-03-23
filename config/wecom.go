package config

import "time"

const (
	WeComEnvKeyCorpID        = "wecom_corp_id"
	WeComEnvKeyCorpSecret    = "wecom_corp_secret"
	WeComEnvKeyRsaPrivateKey = "wecom_rsa_private_key"
)

const (
	WeComChatRecordSDKTimeout = 30 * time.Second
)

type WeComConfig struct {
	CorpID        string
	CorpSecret    string
	RsaPrivateKey string
}

func NewWeComConfig() WeComConfig {
	return WeComConfig{
		CorpID:        lookupEnvVariableFor(WeComEnvKeyCorpID),
		CorpSecret:    lookupEnvVariableFor(WeComEnvKeyCorpSecret),
		RsaPrivateKey: lookupEnvVariableFor(WeComEnvKeyRsaPrivateKey),
	}
}
