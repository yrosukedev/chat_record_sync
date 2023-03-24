package retry_writer

import "github.com/yrosukedev/chat_record_sync/business"

type RetryWriter interface {
	// Write store the chat record.
	//
	// requestUUID is a UUID for retrying the idempotent operation.
	// It's very useful for handling the errors of writing the chat record.
	// Because this operation is idempotent when the UUID is provided, so we could retry the operation many times.
	Write(chatRecord *business.ChatRecord, requestUUID string) error
}
