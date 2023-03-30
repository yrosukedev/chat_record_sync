package paginated_reader

import (
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
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

	inPageToken, err := r.paginationStorage.Get()
	if err != nil {
		return nil, err
	}

	records, outPageToken, err := r.paginatedBufferedReader.Read(inPageToken, r.pageSize)
	if err != nil && err != io.EOF {
		return nil, err
	}

	if uint64(len(records)) < r.pageSize {
		r.noMoreData = true
	}

	if err := r.paginationStorage.Set(outPageToken); err != nil {
		return nil, err
	}

	return records, err
}
