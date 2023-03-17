package paginated_reader

import "github.com/yrosukedev/chat_record_sync/business"

type ChatRecordPaginatedBufferedReader interface {
	Read(pageToken PageToken, pageSize uint64) (records []*business.ChatRecord, err error)
}
