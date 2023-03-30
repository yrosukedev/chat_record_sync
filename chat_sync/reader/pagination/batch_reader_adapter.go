package pagination

import (
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"io"
)

type BatchReaderAdapter struct {
	paginatedReader   PaginatedReader
	paginationStorage PaginationStorage
	pageSize          uint64
	noMoreData        bool
}

func NewBatchReaderAdapter(
	paginatedReader PaginatedReader,
	paginationStorage PaginationStorage,
	pageSize uint64) *BatchReaderAdapter {
	return &BatchReaderAdapter{
		paginatedReader:   paginatedReader,
		paginationStorage: paginationStorage,
		pageSize:          pageSize,
		noMoreData:        false,
	}
}

func (r *BatchReaderAdapter) Read() (records []*business.ChatRecord, err error) {
	if r.noMoreData {
		return nil, io.EOF
	}

	inPageToken, err := r.paginationStorage.Get()
	if err != nil {
		return nil, err
	}

	records, outPageToken, err := r.paginatedReader.Read(inPageToken, r.pageSize)
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
