package paginated_reader

import (
	"github.com/golang/mock/gomock"
	"github.com/yrosukedev/chat_record_sync/business"
	"reflect"
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

func TestForwardResults_zeroRecord(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	bufferedReader := NewMockChatRecordPaginatedBufferedReader(ctrl)
	paginationStorage := NewMockChatRecordPaginationStorage(ctrl)
	paginatedReader := NewChatRecordPaginatedReader(bufferedReader, paginationStorage)

	// Then
	paginationStorage.EXPECT().Get().Return(PageToken(123456), nil).Times(1)

	var records []*business.ChatRecord
	bufferedReader.EXPECT().Read(gomock.Eq(PageToken(123456))).Return(records, nil).Times(1)

	paginationStorage.EXPECT().Set(gomock.Eq(PageToken(123456))).Return(nil).Times(1)

	// When
	forwardingResults, err := paginatedReader.Read()
	if err != nil {
		t.Errorf("error should not happen here, err: %v", err)
	}
	if !reflect.DeepEqual(records, forwardingResults) {
		t.Errorf("the results should not be changed when forwarding it, expetecd: %+v, actual: %+v", records, forwardingResults)
	}
}