package pagination

import (
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
)

type PaginatedReader interface {
	// Read return chat records and the next page token.
	//
	// inPageToken is the input page token. If it's nil, the chat records will be fetched from the beginning.
	Read(inPageToken *PageToken, pageSize uint64) (records []*business.ChatRecord, outPageToken *PageToken, err error)
}
