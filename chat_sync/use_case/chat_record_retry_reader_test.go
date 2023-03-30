package use_case

import (
	"github.com/golang/mock/gomock"
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
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
	GivenRecordsToRead(reader, records)

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
	GivenRecordsToRead(reader, records)

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
	records := []*RecordOrError{
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		NewRecordOrErrorWithError(io.ErrShortBuffer),
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
	}
	EncounterErrorWhileReadingRecords(reader, records)

	// Then
	expectReaderToReadRecordsOrErrors(t, proxyReader, records)
}

func TestOneSequenceOfConsecutiveErrors_oneError_errorCountsEqualToMaxRetryTimes(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	reader := NewMockChatRecordReader(ctrl)
	maxRetryTimes := uint(1)
	proxyReader := NewChatRecordRetryReader(reader, maxRetryTimes)

	// When
	records := []*RecordOrError{
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		NewRecordOrErrorWithError(io.ErrShortBuffer),
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
	}
	EncounterErrorWhileReadingRecords(reader, records)

	// Then
	expectReaderToReadRecordsOrErrors(t, proxyReader, records)
}

func TestOneSequenceOfConsecutiveErrors_oneError_errorCountsGreaterThanMaxRetryTimes(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	reader := NewMockChatRecordReader(ctrl)
	maxRetryTimes := uint(0)
	proxyReader := NewChatRecordRetryReader(reader, maxRetryTimes)

	// When
	records := []*RecordOrError{
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		NewRecordOrErrorWithError(io.ErrShortBuffer),
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
	}
	EncounterErrorWhileReadingRecords(reader, records)

	// Then
	expectedRecordsOrErrors := []*RecordOrError{
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		NewRecordOrErrorWithError(io.ErrShortBuffer),
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
	}
	expectReaderToReadRecordsOrErrors(t, proxyReader, expectedRecordsOrErrors)
}

func TestOneSequenceOfConsecutiveErrors_manyErrors_errorCountsLessThanMaxRetryTimes(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	reader := NewMockChatRecordReader(ctrl)
	maxRetryTimes := uint(4)
	proxyReader := NewChatRecordRetryReader(reader, maxRetryTimes)

	// When
	records := []*RecordOrError{
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		NewRecordOrErrorWithError(io.ErrShortBuffer),
		NewRecordOrErrorWithError(io.ErrNoProgress),
		NewRecordOrErrorWithError(io.ErrClosedPipe),
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
	}
	EncounterErrorWhileReadingRecords(reader, records)

	// Then
	expectReaderToReadRecordsOrErrors(t, proxyReader, records)
}

func TestOneSequenceOfConsecutiveErrors_manyErrors_errorCountsEqualToMaxRetryTimes(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	reader := NewMockChatRecordReader(ctrl)
	maxRetryTimes := uint(4)
	proxyReader := NewChatRecordRetryReader(reader, maxRetryTimes)

	// When
	records := []*RecordOrError{
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		NewRecordOrErrorWithError(io.ErrShortBuffer),
		NewRecordOrErrorWithError(io.ErrUnexpectedEOF),
		NewRecordOrErrorWithError(io.ErrClosedPipe),
		NewRecordOrErrorWithError(io.ErrShortWrite),
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
	}
	EncounterErrorWhileReadingRecords(reader, records)

	// Then
	expectReaderToReadRecordsOrErrors(t, proxyReader, records)
}

func TestOneSequenceOfConsecutiveErrors_manyErrors_errorCountsGreaterThanMaxRetryTimes(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	reader := NewMockChatRecordReader(ctrl)
	maxRetryTimes := uint(3)
	proxyReader := NewChatRecordRetryReader(reader, maxRetryTimes)

	// When
	records := []*RecordOrError{
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		NewRecordOrErrorWithError(io.ErrShortBuffer),
		NewRecordOrErrorWithError(io.ErrUnexpectedEOF),
		NewRecordOrErrorWithError(io.ErrClosedPipe),
		NewRecordOrErrorWithError(io.ErrShortWrite),
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
	}
	EncounterErrorWhileReadingRecords(reader, records)

	// Then
	expectedRecordsOrErrors := []*RecordOrError{
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		NewRecordOrErrorWithError(io.ErrShortBuffer),
		NewRecordOrErrorWithError(io.ErrUnexpectedEOF),
		NewRecordOrErrorWithError(io.ErrClosedPipe),
		NewRecordOrErrorWithError(io.EOF), // terminated from here
		NewRecordOrErrorWithError(io.EOF),
		NewRecordOrErrorWithError(io.EOF),
		NewRecordOrErrorWithError(io.EOF),
	}
	expectReaderToReadRecordsOrErrors(t, proxyReader, expectedRecordsOrErrors)
}

func TestManySequencesOfConsecutiveErrors_errorCountLessThanMaxRetryTimes(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	reader := NewMockChatRecordReader(ctrl)
	maxRetryTimes := uint(5)
	proxyReader := NewChatRecordRetryReader(reader, maxRetryTimes)

	// When
	records := []*RecordOrError{
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		NewRecordOrErrorWithError(io.ErrShortBuffer),
		NewRecordOrErrorWithError(io.ErrUnexpectedEOF),
		NewRecordOrErrorWithError(io.ErrClosedPipe),
		NewRecordOrErrorWithError(io.ErrShortWrite),
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		NewRecordOrErrorWithError(io.ErrNoProgress),
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		NewRecordOrErrorWithError(io.ErrUnexpectedEOF),
		NewRecordOrErrorWithError(io.ErrClosedPipe),
	}
	EncounterErrorWhileReadingRecords(reader, records)

	// Then
	expectReaderToReadRecordsOrErrors(t, proxyReader, records)
}

func TestManySequencesOfConsecutiveErrors_errorCountEqualToMaxRetryTimes(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	reader := NewMockChatRecordReader(ctrl)
	maxRetryTimes := uint(4)
	proxyReader := NewChatRecordRetryReader(reader, maxRetryTimes)

	// When
	records := []*RecordOrError{
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		NewRecordOrErrorWithError(io.ErrUnexpectedEOF),
		NewRecordOrErrorWithError(io.ErrClosedPipe),
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		NewRecordOrErrorWithError(io.ErrNoProgress),
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		NewRecordOrErrorWithError(io.ErrShortBuffer),
		NewRecordOrErrorWithError(io.ErrUnexpectedEOF),
		NewRecordOrErrorWithError(io.ErrClosedPipe),
		NewRecordOrErrorWithError(io.ErrShortWrite),
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
	}
	EncounterErrorWhileReadingRecords(reader, records)

	// Then
	expectReaderToReadRecordsOrErrors(t, proxyReader, records)
}

func TestManySequencesOfConsecutiveErrors_errorCountGreaterThanMaxRetryTimes(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	reader := NewMockChatRecordReader(ctrl)
	maxRetryTimes := uint(3)
	proxyReader := NewChatRecordRetryReader(reader, maxRetryTimes)

	// When
	records := []*RecordOrError{
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		NewRecordOrErrorWithError(io.ErrUnexpectedEOF),
		NewRecordOrErrorWithError(io.ErrClosedPipe),
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		NewRecordOrErrorWithError(io.ErrNoProgress),
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		NewRecordOrErrorWithError(io.ErrShortBuffer),
		NewRecordOrErrorWithError(io.ErrUnexpectedEOF),
		NewRecordOrErrorWithError(io.ErrClosedPipe),
		NewRecordOrErrorWithError(io.ErrShortWrite),
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
	}
	EncounterErrorWhileReadingRecords(reader, records)

	// Then
	expectedRecordsOrErrors := []*RecordOrError{
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		NewRecordOrErrorWithError(io.ErrUnexpectedEOF),
		NewRecordOrErrorWithError(io.ErrClosedPipe),
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		NewRecordOrErrorWithError(io.ErrNoProgress),
		NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		NewRecordOrErrorWithError(io.ErrShortBuffer),
		NewRecordOrErrorWithError(io.ErrUnexpectedEOF),
		NewRecordOrErrorWithError(io.ErrClosedPipe),
		NewRecordOrErrorWithError(io.EOF),
	}
	expectReaderToReadRecordsOrErrors(t, proxyReader, expectedRecordsOrErrors)
}

func expectReaderToReadRecordsOrErrors(t *testing.T, reader ChatRecordReader, records []*RecordOrError) {
	for _, record := range records {
		switch record.InnerType {
		case RecordOrErrorInnerTypeError:
			if _, err := reader.Read(); err != record.Err {
				t.Errorf("error should happen here, expected: %+v, actual: %+v", record.Err, err)
			}
		case RecordOrErrorInnerTypeRecord:
			if _, err := reader.Read(); err != nil {
				t.Errorf("error should not happen here, expected: %+v, actual: %+v", nil, err)
			}
		}
	}

	if _, err := reader.Read(); err != io.EOF {
		t.Errorf("io.EOF should be returned here, expected: %+v, actual: %+v", io.EOF, err)
	}
}
