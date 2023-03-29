package cli_app

import (
	"context"
	"errors"
	lark "github.com/larksuite/oapi-sdk-go/v3"
	"github.com/yrosukedev/WeWorkFinanceSDK"
	"github.com/yrosukedev/chat_record_sync/buffer_reader"
	"github.com/yrosukedev/chat_record_sync/chat_record_bitable_storage"
	"github.com/yrosukedev/chat_record_sync/config"
	logproxy "github.com/yrosukedev/chat_record_sync/logger/proxy"
	"github.com/yrosukedev/chat_record_sync/paginated_reader"
	"github.com/yrosukedev/chat_record_sync/pagination_bitable_storage"
	"github.com/yrosukedev/chat_record_sync/retry_writer"
	"github.com/yrosukedev/chat_record_sync/use_case"
	"github.com/yrosukedev/chat_record_sync/wecom_chat"
	"github.com/yrosukedev/chat_record_sync/wecom_chat_adapter"
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

	useCase := use_case.NewSyncChatRecordUseCase(
		buffer_reader.NewChatRecordBufferedReaderAdapter(
			paginated_reader.NewChatRecordPaginatedReader(
				wecom_chat.NewPaginatedBufferedReaderAdapter(
					wecom_chat_adapter.NewWeComChatRecordServiceAdapter(ctx, client, "", "", config.WeComChatRecordSDKTimeout, logger),
					nil,
					wecom_chat.NewWeComMessageTransformerFactory(map[string]wecom_chat.ChatRecordTransformer{
						wecom_chat.WeComMessageTypeText: wecom_chat.NewWeComTextMessageTransformer(ctx, logger),
					},
						wecom_chat.NewWeComDefaultMessageTransformer(ctx, logger))),
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
