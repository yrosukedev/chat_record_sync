package use_case

import (
	"github.com/golang/mock/gomock"
	"github.com/yrosukedev/chat_record_sync/business"
	"io"
	"testing"
)

func TestZeroSequenceOfConsecutiveErrors_zeroRead(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	reader := NewMockChatRecordReader(ctrl)
	maxRetryTimes := uint(3)
	proxyReader := NewChatRecordRetryReader(reader, maxRetryTimes)

	// When
	reader.EXPECT().Read().Return(nil, io.EOF).Times(1)

	// Then
	_, err := proxyReader.Read()
	if err != io.EOF {
		t.Error("io.EOF should be returned here")
	}
}

func TestZeroSequenceOfConsecutiveErrors_oneRead(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	reader := NewMockChatRecordReader(ctrl)
	maxRetryTimes := uint(3)
	proxyReader := NewChatRecordRetryReader(reader, maxRetryTimes)

	// When
	records := []*business.ChatRecord{
		{},
	}
	givenRecordsToRead(reader, records)

	// Then
	if _, err := proxyReader.Read(); err != nil {
		t.Error("error should not happen here")
	}

	_, err := proxyReader.Read()
	if err != io.EOF {
		t.Error("io.EOF should be returned here")
	}
}

func TestZeroSequenceOfConsecutiveErrors_manyReads(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	reader := NewMockChatRecordReader(ctrl)
	maxRetryTimes := uint(3)
	proxyReader := NewChatRecordRetryReader(reader, maxRetryTimes)

	// When
	records := []*business.ChatRecord{
		{},
		{},
		{},
	}
	givenRecordsToRead(reader, records)

	// Then
	for i := 0; i < len(records); i++ {
		if _, err := proxyReader.Read(); err != nil {
			t.Error("error should not happen here")
		}
	}

	_, err := proxyReader.Read()
	if err != io.EOF {
		t.Error("io.EOF should be returned here")
	}
}
