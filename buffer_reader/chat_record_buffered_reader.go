package buffer_reader

import "github.com/yrosukedev/chat_record_sync/business"

type ChatRecordBufferedReader interface {
	Read() (records []*business.ChatRecord, err error)
}
