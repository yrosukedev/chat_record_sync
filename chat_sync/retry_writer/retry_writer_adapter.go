package retry_writer

import (
	"github.com/google/uuid"
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/chat_sync/use_case"
)

type RetryWriterAdapter struct {
	retryWriter RetryWriter
}

func NewRetryWriterAdapter(retryWriter RetryWriter) use_case.Writer {
	return &RetryWriterAdapter{
		retryWriter: retryWriter,
	}
}

func (r *RetryWriterAdapter) Write(record *business.ChatRecord) error {
	// TODO: retry many times
	return r.retryWriter.Write(record, uuid.NewString())
}
