package chat_record

import "github.com/yrosukedev/chat_record_sync/chat_sync/business"

type FieldsWriter struct {
	fieldsFormatter FieldsFormatter
	fieldsStorage   FieldsStorage
}

func NewFieldsWriter(fieldsFormatter FieldsFormatter, fieldsStorage FieldsStorage) *FieldsWriter {
	return &FieldsWriter{
		fieldsFormatter: fieldsFormatter,
		fieldsStorage:   fieldsStorage,
	}
}

func (w *FieldsWriter) Write(chatRecord *business.ChatRecord, requestUUID string) error {
	fields, err := w.fieldsFormatter.Format(chatRecord)
	if err != nil {
		return err
	}

	if len(fields) == 0 {
		return nil
	}

	return w.fieldsStorage.Write(fields, requestUUID)
}
