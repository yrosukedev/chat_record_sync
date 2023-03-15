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
	if _, err := proxyReader.Read(); err != io.EOF {
		t.Errorf("io.EOF should be returned here, expected: %+v, actual: %+v", io.EOF, err)
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
		t.Errorf("error should not happen here, expected: %+v, actual: %+v", nil, err)
	}

	if _, err := proxyReader.Read(); err != io.EOF {
		t.Errorf("io.EOF should be returned here, expected: %+v, actual: %+v", io.EOF, err)
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
			t.Errorf("error should not happen here, expected: %+v, actual: %+v", nil, err)
		}
	}

	if _, err := proxyReader.Read(); err != io.EOF {
		t.Errorf("io.EOF should be returned here, expected: %+v, actual: %+v", io.EOF, err)
	}
}

func TestOneSequenceOfConsecutiveErrors_oneError_errorCountsLessThanMaxRetryTimes(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	reader := NewMockChatRecordReader(ctrl)
	maxRetryTimes := uint(3)
	proxyReader := NewChatRecordRetryReader(reader, maxRetryTimes)

	// When
	records := []*recordOrError{
		newRecordOrErrorWithRecord(&business.ChatRecord{}),
		newRecordOrErrorWithError(io.ErrShortBuffer),
		newRecordOrErrorWithRecord(&business.ChatRecord{}),
	}
	encounterErrorWhileReadingRecords(reader, records)

	// Then
	if _, err := proxyReader.Read(); err != nil {
		t.Errorf("error should not happen here, expected: %+v, actual: %+v", nil, err)
	}
	if _, err := proxyReader.Read(); err != io.ErrShortBuffer {
		t.Errorf("error should happen here, expected: %+v, actual: %+v", io.ErrShortBuffer, err)
	}
	if _, err := proxyReader.Read(); err != nil {
		t.Errorf("error should not happen here, expected: %+v, actual: %+v", nil, err)
	}
	if _, err := proxyReader.Read(); err != io.EOF {
		t.Errorf("io.EOF should be returned here, expected: %+v, actual: %+v", io.EOF, err)
	}
}

func TestOneSequenceOfConsecutiveErrors_oneError_errorCountsEqualToMaxRetryTimes(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	reader := NewMockChatRecordReader(ctrl)
	maxRetryTimes := uint(1)
	proxyReader := NewChatRecordRetryReader(reader, maxRetryTimes)

	// When
	records := []*recordOrError{
		newRecordOrErrorWithRecord(&business.ChatRecord{}),
		newRecordOrErrorWithError(io.ErrShortBuffer),
		newRecordOrErrorWithRecord(&business.ChatRecord{}),
	}
	encounterErrorWhileReadingRecords(reader, records)

	// Then
	if _, err := proxyReader.Read(); err != nil {
		t.Errorf("error should not happen here, expected: %+v, actual: %+v", nil, err)
	}
	if _, err := proxyReader.Read(); err != io.ErrShortBuffer {
		t.Errorf("error should happen here, expected: %+v, actual: %+v", io.ErrShortBuffer, err)
	}
	if _, err := proxyReader.Read(); err != nil {
		t.Errorf("error should not happen here, expected: %+v, actual: %+v", nil, err)
	}
	if _, err := proxyReader.Read(); err != io.EOF {
		t.Errorf("io.EOF should be returned here, expected: %+v, actual: %+v", io.EOF, err)
	}
}
