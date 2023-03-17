package paginated_reader

import (
	"github.com/yrosukedev/chat_record_sync/business"
	"io"
)

type ChatRecordPaginatedReader struct {
	paginatedBufferedReader ChatRecordPaginatedBufferedReader
	paginationStorage       ChatRecordPaginationStorage
	pageSize                uint64
	noMoreData              bool
}

func NewChatRecordPaginatedReader(
	paginatedBufferedReader ChatRecordPaginatedBufferedReader,
	paginationStorage ChatRecordPaginationStorage,
	pageSize uint64) *ChatRecordPaginatedReader {
	return &ChatRecordPaginatedReader{
		paginatedBufferedReader: paginatedBufferedReader,
		paginationStorage:       paginationStorage,
		pageSize:                pageSize,
		noMoreData:              false,
	}
}

func (r *ChatRecordPaginatedReader) Read() (records []*business.ChatRecord, err error) {
	if r.noMoreData {
		return nil, io.EOF
	}

	pageToken, err := r.paginationStorage.Get()
	if err != nil {
		return nil, err
	}

	records, err = r.paginatedBufferedReader.Read(pageToken, r.pageSize)
	if err != nil {
		return nil, err
	}

	if uint64(len(records)) < r.pageSize {
		r.noMoreData = true
	}

	if err := r.paginationStorage.Set(pageToken + PageToken(len(records))); err != nil {
		return nil, err
	}

	return records, err
}
