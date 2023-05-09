package chat_record

import (
	"context"
	"errors"
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/logger"
)

type FieldsWriter struct {
	ctx             context.Context
	fieldsFormatter FieldsFormatter
	fieldsStorage   FieldsStorage
	logger          logger.Logger
}

func NewFieldsWriter(ctx context.Context, fieldsFormatter FieldsFormatter, fieldsStorage FieldsStorage, logger logger.Logger) *FieldsWriter {
	return &FieldsWriter{
		ctx:             ctx,
		fieldsFormatter: fieldsFormatter,
		fieldsStorage:   fieldsStorage,
		logger:          logger,
	}
}

func (w *FieldsWriter) Write(chatRecord *business.ChatRecord, requestUUID string) error {
	if chatRecord == nil {
		w.logger.Error(w.ctx, "[chat record fields writer] chatRecord can't be nil")
		return errors.New("chatRecord can't be nil")
	}

	w.logger.Info(w.ctx, "[chat record fields writer] will write record to storage, msgId: %v", chatRecord.MsgId)

	fields, err := w.fieldsFormatter.Format(chatRecord)
	if err != nil {
		w.logger.Error(w.ctx, "[chat record fields writer] fails to format fields, msgId: %v, err: %v", chatRecord.MsgId, err)
		return err
	}

	if len(fields) == 0 {
		w.logger.Info(w.ctx, "[chat record fields writer] fields is empty, won't write them to storage, msgId: %v", chatRecord.MsgId)
		return nil
	}

	err = w.fieldsStorage.Write(fields, requestUUID)
	if err != nil {
		w.logger.Error(w.ctx, "[chat record fields writer] fails to write fields to storage, msgId: %v, err: %v", chatRecord.MsgId, err)
		return err
	}

	w.logger.Info(w.ctx, "[chat record fields writer] successfully write record to storage, msgId: %v", chatRecord.MsgId)

	return nil

}
