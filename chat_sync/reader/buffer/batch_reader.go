package buffer

import (
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
)

type BatchReader interface {
	Read() (records []*business.ChatRecord, err error)
}
