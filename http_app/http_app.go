package http_app

import (
	"context"
	"fmt"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/xen0n/go-workwx"
	"github.com/yrosukedev/WeWorkFinanceSDK"
	"github.com/yrosukedev/chat_record_sync/config"
	logproxy "github.com/yrosukedev/chat_record_sync/logger/proxy"
	"net/http"

	log "github.com/yrosukedev/chat_record_sync/logger"
)

type HTTPApp struct {
	larkClient  *lark.Client
	wecomConfig config.WeComConfig
	wecomClient WeWorkFinanceSDK.Client
	wecomApp    *workwx.WorkwxApp
	logger      log.Logger
}

func NewHTTPApp(ctx context.Context) *HTTPApp {
	logger := logproxy.NewLoggerProxy(config.HttpAppLogLevel, logproxy.NewDefaultLogger())

	logger.Info(ctx, "[http app] start to create http app")

	larkConfig := config.NewLarkConfig()
	larkClient := lark.NewClient(larkConfig.AppId, larkConfig.AppSecret)
	logger.Info(ctx, "[http app] lark client created")

	weComConfig := config.NewWeComConfig()
	wecomClient, err := WeWorkFinanceSDK.NewClient(weComConfig.CorpID, weComConfig.ChatSyncSecret, weComConfig.ChatSyncRsaPrivateKey)
	if err != nil {
		logger.Error(ctx, "[http app] fails to create wecom client, err: %v", err)
		panic(fmt.Sprintf("fails to create wecom client, err: %v", err))
	}
	logger.Info(ctx, "[http app] wecom client created")

	wecomApp := workwx.New(weComConfig.CorpID).WithApp(weComConfig.AgentSecret, weComConfig.AgentID)

	httpApp := &HTTPApp{
		larkClient:  larkClient,
		wecomConfig: weComConfig,
		wecomClient: wecomClient,
		wecomApp:    wecomApp,
		logger:      logger,
	}

	logger.Info(ctx, "[http app] http app created")

	return httpApp
}

func (f *HTTPApp) Run(ctx context.Context) {
	f.logger.Info(ctx, "[http app] start to run server")

	requestMux := f.createMultiplexer(ctx)
	f.logger.Info(ctx, "[http app] request multiplexer created")

	err := http.ListenAndServe(
		fmt.Sprintf(":%v", config.HttpAppPort),
		NewLogHandler(ctx, requestMux, f.logger))
	if err != nil {
		f.logger.Error(ctx, "[http app] server exit with err: %v", err)
	} else {
		f.logger.Info(ctx, "[http app] server stopped")
	}
}
