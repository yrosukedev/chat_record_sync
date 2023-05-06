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
	wecom_chat2 "github.com/yrosukedev/chat_record_sync/chat_sync/wecom"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom/chat_record_service"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom/openapi"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom/transformer"
	"github.com/yrosukedev/chat_record_sync/chat_sync/writer/retry"
	"github.com/yrosukedev/chat_record_sync/config"
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

	wecomApp := workwx.New(weComConfig.CorpID).WithApp(weComConfig.AgentSecret, weComConfig.AgentID)

	logger := logproxy.NewLoggerProxy(config.HttpAppLogLevel, logproxy.NewDefaultLogger())

	useCase := use_case.NewChatSyncUseCase(
		buffer.NewReader(
			pagination.NewBatchReaderAdapter(
				wecom_chat2.NewPaginatedReaderAdapter(
					chat_record_service.NewAdapter(ctx, client, "", "", config.WeComChatRecordSDKTimeout, logger),
					transformer.NewRecordTransformerBuilder(
						openapi.NewAdapter(ctx, wecomApp, logger)).Build()),
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
