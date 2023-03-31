package cli_app

import (
	"context"
	"errors"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/yrosukedev/WeWorkFinanceSDK"
	"github.com/yrosukedev/chat_record_sync/chat_sync/chat_record_bitable_storage"
	"github.com/yrosukedev/chat_record_sync/chat_sync/pagination_bitable_storage"
	"github.com/yrosukedev/chat_record_sync/chat_sync/reader/buffer"
	"github.com/yrosukedev/chat_record_sync/chat_sync/reader/pagination"
	"github.com/yrosukedev/chat_record_sync/chat_sync/retry_writer"
	"github.com/yrosukedev/chat_record_sync/chat_sync/use_case"
	wecom_chat2 "github.com/yrosukedev/chat_record_sync/chat_sync/wecom"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom/transformer"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom_chat_adapter"
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

	pageSize := uint64(10)

	logger := logproxy.NewLoggerProxy(config.HttpAppLogLevel, logproxy.NewDefaultLogger())

	useCase := use_case.NewChatSyncUseCase(
		buffer.NewReader(
			pagination.NewBatchReaderAdapter(
				wecom_chat2.NewPaginatedReaderAdapter(
					wecom_chat_adapter.NewWeComChatRecordServiceAdapter(ctx, client, "", "", config.WeComChatRecordSDKTimeout, logger),
					nil,
					transformer.NewWeComMessageTransformerFactory(map[string]wecom_chat2.ChatRecordTransformer{
						wecom_chat2.MessageTypeText: transformer.NewWeComTextMessageTransformer(ctx, logger),
					},
						transformer.NewWeComDefaultMessageTransformer(ctx, logger))),
				pagination_bitable_storage.NewPaginationStorageAdapter(ctx, larkClient, "DLSbbQIcEa0KyIsetHWcg3PDnNh", "tblLJY5YSoEkV3G3", logger),
				pageSize)),
		retry_writer.NewRetryWriterAdapter(
			chat_record_bitable_storage.NewChatRecordStorageAdapter(ctx, larkClient, "QCBrbzgx4aKRAis9eewcV731n7d", "tblIk692K5LXte8x", logger)),
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
