package buffer_reader

import (
	"github.com/golang/mock/gomock"
	"github.com/yrosukedev/chat_record_sync/business"
	"testing"
)

func TestBufferSize_one(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	bufferedReader := NewMockChatRecordBufferedReader(ctrl)
	readerAdapter := NewChatRecordBufferedReaderAdapter(bufferedReader)

	// Then
	records := []*business.ChatRecord{
		{},
	}
	bufferedReader.EXPECT().Read().Return(records, nil).Times(10)

	// When
	for i := 0; i < 10; i++ {
		_, _ = readerAdapter.Read()
	}
}

func TestBufferSize_zero(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	bufferedReader := NewMockChatRecordBufferedReader(ctrl)
	readerAdapter := NewChatRecordBufferedReaderAdapter(bufferedReader)

	// Then
	var records []*business.ChatRecord
	bufferedReader.EXPECT().Read().Return(records, nil).Times(10)

	// When
	for i := 0; i < 10; i++ {
		_, _ = readerAdapter.Read()
	}
}
