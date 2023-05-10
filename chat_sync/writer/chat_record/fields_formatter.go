package chat_record

import "github.com/yrosukedev/chat_record_sync/chat_sync/business"

type FieldsFormatter interface {
	Format(record *business.ChatRecord) (fields map[string]interface{}, err error)
}
