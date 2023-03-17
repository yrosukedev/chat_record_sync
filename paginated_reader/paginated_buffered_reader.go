package paginated_reader

import "github.com/yrosukedev/chat_record_sync/business"

type ChatRecordPaginatedBufferedReader interface {
	Read(pageToken PageToken) (records []*business.ChatRecord, err error)
}
