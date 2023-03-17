package paginated_reader

import (
	"github.com/golang/mock/gomock"
	"github.com/yrosukedev/chat_record_sync/business"
	"testing"
)

func TestFetchPageToken_zero(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	bufferedReader := NewMockChatRecordPaginatedBufferedReader(ctrl)
	paginationStorage := NewMockChatRecordPaginationStorage(ctrl)
	paginatedReader := NewChatRecordPaginatedReader(bufferedReader, paginationStorage)

	// Then
	paginationStorage.EXPECT().Get().Return(PageToken(0), nil).Times(1)
	bufferedReader.EXPECT().Read(gomock.Eq(PageToken(0))).Return([]*business.ChatRecord{}, nil).Times(1)

	// When
	paginatedReader.Read()
}

func TestFetchPageToken_one(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	bufferedReader := NewMockChatRecordPaginatedBufferedReader(ctrl)
	paginationStorage := NewMockChatRecordPaginationStorage(ctrl)
	paginatedReader := NewChatRecordPaginatedReader(bufferedReader, paginationStorage)

	// Then
	paginationStorage.EXPECT().Get().Return(PageToken(1), nil).Times(1)
	bufferedReader.EXPECT().Read(gomock.Eq(PageToken(1))).Return([]*business.ChatRecord{}, nil).Times(1)

	// When
	paginatedReader.Read()
}

func TestFetchPageToken_many(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	bufferedReader := NewMockChatRecordPaginatedBufferedReader(ctrl)
	paginationStorage := NewMockChatRecordPaginationStorage(ctrl)
	paginatedReader := NewChatRecordPaginatedReader(bufferedReader, paginationStorage)

	// Then
	paginationStorage.EXPECT().Get().Return(PageToken(3), nil).Times(1)
	bufferedReader.EXPECT().Read(gomock.Eq(PageToken(3))).Return([]*business.ChatRecord{}, nil).Times(1)

	// When
	paginatedReader.Read()
}

func TestFetchPageToken_ALot(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	bufferedReader := NewMockChatRecordPaginatedBufferedReader(ctrl)
	paginationStorage := NewMockChatRecordPaginationStorage(ctrl)
	paginatedReader := NewChatRecordPaginatedReader(bufferedReader, paginationStorage)

	// Then
	paginationStorage.EXPECT().Get().Return(PageToken(10500), nil).Times(1)
	bufferedReader.EXPECT().Read(gomock.Eq(PageToken(10500))).Return([]*business.ChatRecord{}, nil).Times(1)

	// When
	paginatedReader.Read()
}
