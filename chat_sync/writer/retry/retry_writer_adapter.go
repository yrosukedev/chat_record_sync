package retry

import (
	"github.com/google/uuid"
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/chat_sync/use_case"
)

type WriterAdapter struct {
	retryWriter RetryWriter
}

func NewWriterAdapter(retryWriter RetryWriter) use_case.Writer {
	return &WriterAdapter{
		retryWriter: retryWriter,
	}
}

func (r *WriterAdapter) Write(record *business.ChatRecord) error {
	// TODO: retry many times
	return r.retryWriter.Write(record, uuid.NewString())
}
