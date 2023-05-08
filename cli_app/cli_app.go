package cli_app

import (
	"context"
	"errors"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/xen0n/go-workwx"
	"github.com/yrosukedev/WeWorkFinanceSDK"
	"github.com/yrosukedev/chat_record_sync/chat_sync/bitable_storage/chat_record"
	pagination2 "github.com/yrosukedev/chat_record_sync/chat_sync/bitable_storage/pagination"
	"github.com/yrosukedev/chat_record_sync/chat_sync/reader/buffer"
	"github.com/yrosukedev/chat_record_sync/chat_sync/reader/pagination"
	"github.com/yrosukedev/chat_record_sync/chat_sync/use_case"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom/chat_record_service"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom/openapi"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom/transformer"
	"github.com/yrosukedev/chat_record_sync/chat_sync/writer/retry"
	"github.com/yrosukedev/chat_record_sync/config"
	"github.com/yrosukedev/chat_record_sync/logger"
	logproxy "github.com/yrosukedev/chat_record_sync/logger/proxy"
	"strings"
)

func RunCLIApp(ctx context.Context) error {
	larkConfig := config.NewLarkConfig()
	larkClient := lark.NewClient(larkConfig.AppId, larkConfig.AppSecret)

	weComConfig := config.NewWeComConfig()
	client, err := WeWorkFinanceSDK.NewClient(weComConfig.CorpID, weComConfig.ChatSyncSecret, weComConfig.ChatSyncRsaPrivateKey)
	if err != nil {
		return err
	}

	logger := logproxy.NewLoggerProxy(config.HttpAppLogLevel, logproxy.NewDefaultLogger())

	useCase := use_case.NewChatSyncUseCase(
		buffer.NewReader(
			pagination.NewBatchReaderAdapter(
				wecom.NewPaginatedReaderAdapter(
					chat_record_service.NewAdapter(ctx, client, "", "", config.WeComChatRecordSDKTimeout, logger),
					buildRecordTransformer(ctx, weComConfig, logger)),
				pagination2.NewStorageAdapter(ctx, larkClient, config.PaginationStorageBitableAppToken, config.PaginationStorageBitableTableId, logger),
				config.PaginatedReaderPageSize)),
		retry.NewWriterAdapter(
			chat_record.NewStorageAdapter(ctx, larkClient, config.ChatStorageBitableAppToken, config.ChatStorageBitableTableId, logger)),
	)

	// When
	errs := useCase.Run(ctx)

	// Then
	if len(errs) != 0 {
		var errMsgs []string
		for _, err := range errs {
			errMsgs = append(errMsgs, err.Error())
		}
		return errors.New(strings.Join(errMsgs, "\n"))
	}

	return nil
}

func buildRecordTransformer(ctx context.Context, weComConfig config.WeComConfig, logger logger.Logger) wecom.RecordTransformer {
	wecomApp := workwx.New(weComConfig.CorpID).WithApp(weComConfig.AgentSecret, weComConfig.AgentID)
	msgAuditWecomApp := workwx.New(weComConfig.CorpID).WithApp(weComConfig.ChatSyncSecret, config.WeComMsgAuditAgentID)
	return transformer.NewRecordTransformerBuilder(
		openapi.NewAdapter(ctx, wecomApp, logger),
		openapi.NewMsgAuditOpenAPIAdapter(ctx, msgAuditWecomApp, logger)).Build()
}
