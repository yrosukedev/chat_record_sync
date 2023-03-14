package use_case

import "github.com/yrosukedev/chat_record_sync/business"

type ChatRecordWriter interface {
	Write(record *business.ChatRecord) error
}
