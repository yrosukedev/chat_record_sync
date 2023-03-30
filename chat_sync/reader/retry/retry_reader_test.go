package retry

import (
	"github.com/golang/mock/gomock"
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/chat_sync/use_case"
	"io"
	"testing"
)

func TestZeroSequenceOfConsecutiveErrors_zeroRead(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	reader := use_case.NewMockReader(ctrl)
	maxRetryTimes := uint(3)
	proxyReader := NewRetryReader(reader, maxRetryTimes)

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
	reader := use_case.NewMockReader(ctrl)
	maxRetryTimes := uint(3)
	proxyReader := NewRetryReader(reader, maxRetryTimes)

	// When
	records := []*business.ChatRecord{
		{},
	}
	use_case.GivenRecordsToRead(reader, records)

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
	reader := use_case.NewMockReader(ctrl)
	maxRetryTimes := uint(3)
	proxyReader := NewRetryReader(reader, maxRetryTimes)

	// When
	records := []*business.ChatRecord{
		{},
		{},
		{},
	}
	use_case.GivenRecordsToRead(reader, records)

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
	reader := use_case.NewMockReader(ctrl)
	maxRetryTimes := uint(3)
	proxyReader := NewRetryReader(reader, maxRetryTimes)

	// When
	records := []*use_case.RecordOrError{
		use_case.NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		use_case.NewRecordOrErrorWithError(io.ErrShortBuffer),
		use_case.NewRecordOrErrorWithRecord(&business.ChatRecord{}),
	}
	use_case.EncounterErrorWhileReadingRecords(reader, records)

	// Then
	expectReaderToReadRecordsOrErrors(t, proxyReader, records)
}

func TestOneSequenceOfConsecutiveErrors_oneError_errorCountsEqualToMaxRetryTimes(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	reader := use_case.NewMockReader(ctrl)
	maxRetryTimes := uint(1)
	proxyReader := NewRetryReader(reader, maxRetryTimes)

	// When
	records := []*use_case.RecordOrError{
		use_case.NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		use_case.NewRecordOrErrorWithError(io.ErrShortBuffer),
		use_case.NewRecordOrErrorWithRecord(&business.ChatRecord{}),
	}
	use_case.EncounterErrorWhileReadingRecords(reader, records)

	// Then
	expectReaderToReadRecordsOrErrors(t, proxyReader, records)
}

func TestOneSequenceOfConsecutiveErrors_oneError_errorCountsGreaterThanMaxRetryTimes(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	reader := use_case.NewMockReader(ctrl)
	maxRetryTimes := uint(0)
	proxyReader := NewRetryReader(reader, maxRetryTimes)

	// When
	records := []*use_case.RecordOrError{
		use_case.NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		use_case.NewRecordOrErrorWithError(io.ErrShortBuffer),
		use_case.NewRecordOrErrorWithRecord(&business.ChatRecord{}),
	}
	use_case.EncounterErrorWhileReadingRecords(reader, records)

	// Then
	expectedRecordsOrErrors := []*use_case.RecordOrError{
		use_case.NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		use_case.NewRecordOrErrorWithError(io.ErrShortBuffer),
		use_case.NewRecordOrErrorWithRecord(&business.ChatRecord{}),
	}
	expectReaderToReadRecordsOrErrors(t, proxyReader, expectedRecordsOrErrors)
}

func TestOneSequenceOfConsecutiveErrors_manyErrors_errorCountsLessThanMaxRetryTimes(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	reader := use_case.NewMockReader(ctrl)
	maxRetryTimes := uint(4)
	proxyReader := NewRetryReader(reader, maxRetryTimes)

	// When
	records := []*use_case.RecordOrError{
		use_case.NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		use_case.NewRecordOrErrorWithError(io.ErrShortBuffer),
		use_case.NewRecordOrErrorWithError(io.ErrNoProgress),
		use_case.NewRecordOrErrorWithError(io.ErrClosedPipe),
		use_case.NewRecordOrErrorWithRecord(&business.ChatRecord{}),
	}
	use_case.EncounterErrorWhileReadingRecords(reader, records)

	// Then
	expectReaderToReadRecordsOrErrors(t, proxyReader, records)
}

func TestOneSequenceOfConsecutiveErrors_manyErrors_errorCountsEqualToMaxRetryTimes(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	reader := use_case.NewMockReader(ctrl)
	maxRetryTimes := uint(4)
	proxyReader := NewRetryReader(reader, maxRetryTimes)

	// When
	records := []*use_case.RecordOrError{
		use_case.NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		use_case.NewRecordOrErrorWithError(io.ErrShortBuffer),
		use_case.NewRecordOrErrorWithError(io.ErrUnexpectedEOF),
		use_case.NewRecordOrErrorWithError(io.ErrClosedPipe),
		use_case.NewRecordOrErrorWithError(io.ErrShortWrite),
		use_case.NewRecordOrErrorWithRecord(&business.ChatRecord{}),
	}
	use_case.EncounterErrorWhileReadingRecords(reader, records)

	// Then
	expectReaderToReadRecordsOrErrors(t, proxyReader, records)
}

func TestOneSequenceOfConsecutiveErrors_manyErrors_errorCountsGreaterThanMaxRetryTimes(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	reader := use_case.NewMockReader(ctrl)
	maxRetryTimes := uint(3)
	proxyReader := NewRetryReader(reader, maxRetryTimes)

	// When
	records := []*use_case.RecordOrError{
		use_case.NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		use_case.NewRecordOrErrorWithError(io.ErrShortBuffer),
		use_case.NewRecordOrErrorWithError(io.ErrUnexpectedEOF),
		use_case.NewRecordOrErrorWithError(io.ErrClosedPipe),
		use_case.NewRecordOrErrorWithError(io.ErrShortWrite),
		use_case.NewRecordOrErrorWithRecord(&business.ChatRecord{}),
	}
	use_case.EncounterErrorWhileReadingRecords(reader, records)

	// Then
	expectedRecordsOrErrors := []*use_case.RecordOrError{
		use_case.NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		use_case.NewRecordOrErrorWithError(io.ErrShortBuffer),
		use_case.NewRecordOrErrorWithError(io.ErrUnexpectedEOF),
		use_case.NewRecordOrErrorWithError(io.ErrClosedPipe),
		use_case.NewRecordOrErrorWithError(io.EOF), // terminated from here
		use_case.NewRecordOrErrorWithError(io.EOF),
		use_case.NewRecordOrErrorWithError(io.EOF),
		use_case.NewRecordOrErrorWithError(io.EOF),
	}
	expectReaderToReadRecordsOrErrors(t, proxyReader, expectedRecordsOrErrors)
}

func TestManySequencesOfConsecutiveErrors_errorCountLessThanMaxRetryTimes(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	reader := use_case.NewMockReader(ctrl)
	maxRetryTimes := uint(5)
	proxyReader := NewRetryReader(reader, maxRetryTimes)

	// When
	records := []*use_case.RecordOrError{
		use_case.NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		use_case.NewRecordOrErrorWithError(io.ErrShortBuffer),
		use_case.NewRecordOrErrorWithError(io.ErrUnexpectedEOF),
		use_case.NewRecordOrErrorWithError(io.ErrClosedPipe),
		use_case.NewRecordOrErrorWithError(io.ErrShortWrite),
		use_case.NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		use_case.NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		use_case.NewRecordOrErrorWithError(io.ErrNoProgress),
		use_case.NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		use_case.NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		use_case.NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		use_case.NewRecordOrErrorWithError(io.ErrUnexpectedEOF),
		use_case.NewRecordOrErrorWithError(io.ErrClosedPipe),
	}
	use_case.EncounterErrorWhileReadingRecords(reader, records)

	// Then
	expectReaderToReadRecordsOrErrors(t, proxyReader, records)
}

func TestManySequencesOfConsecutiveErrors_errorCountEqualToMaxRetryTimes(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	reader := use_case.NewMockReader(ctrl)
	maxRetryTimes := uint(4)
	proxyReader := NewRetryReader(reader, maxRetryTimes)

	// When
	records := []*use_case.RecordOrError{
		use_case.NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		use_case.NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		use_case.NewRecordOrErrorWithError(io.ErrUnexpectedEOF),
		use_case.NewRecordOrErrorWithError(io.ErrClosedPipe),
		use_case.NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		use_case.NewRecordOrErrorWithError(io.ErrNoProgress),
		use_case.NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		use_case.NewRecordOrErrorWithError(io.ErrShortBuffer),
		use_case.NewRecordOrErrorWithError(io.ErrUnexpectedEOF),
		use_case.NewRecordOrErrorWithError(io.ErrClosedPipe),
		use_case.NewRecordOrErrorWithError(io.ErrShortWrite),
		use_case.NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		use_case.NewRecordOrErrorWithRecord(&business.ChatRecord{}),
	}
	use_case.EncounterErrorWhileReadingRecords(reader, records)

	// Then
	expectReaderToReadRecordsOrErrors(t, proxyReader, records)
}

func TestManySequencesOfConsecutiveErrors_errorCountGreaterThanMaxRetryTimes(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	reader := use_case.NewMockReader(ctrl)
	maxRetryTimes := uint(3)
	proxyReader := NewRetryReader(reader, maxRetryTimes)

	// When
	records := []*use_case.RecordOrError{
		use_case.NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		use_case.NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		use_case.NewRecordOrErrorWithError(io.ErrUnexpectedEOF),
		use_case.NewRecordOrErrorWithError(io.ErrClosedPipe),
		use_case.NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		use_case.NewRecordOrErrorWithError(io.ErrNoProgress),
		use_case.NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		use_case.NewRecordOrErrorWithError(io.ErrShortBuffer),
		use_case.NewRecordOrErrorWithError(io.ErrUnexpectedEOF),
		use_case.NewRecordOrErrorWithError(io.ErrClosedPipe),
		use_case.NewRecordOrErrorWithError(io.ErrShortWrite),
		use_case.NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		use_case.NewRecordOrErrorWithRecord(&business.ChatRecord{}),
	}
	use_case.EncounterErrorWhileReadingRecords(reader, records)

	// Then
	expectedRecordsOrErrors := []*use_case.RecordOrError{
		use_case.NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		use_case.NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		use_case.NewRecordOrErrorWithError(io.ErrUnexpectedEOF),
		use_case.NewRecordOrErrorWithError(io.ErrClosedPipe),
		use_case.NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		use_case.NewRecordOrErrorWithError(io.ErrNoProgress),
		use_case.NewRecordOrErrorWithRecord(&business.ChatRecord{}),
		use_case.NewRecordOrErrorWithError(io.ErrShortBuffer),
		use_case.NewRecordOrErrorWithError(io.ErrUnexpectedEOF),
		use_case.NewRecordOrErrorWithError(io.ErrClosedPipe),
		use_case.NewRecordOrErrorWithError(io.EOF),
	}
	expectReaderToReadRecordsOrErrors(t, proxyReader, expectedRecordsOrErrors)
}

func expectReaderToReadRecordsOrErrors(t *testing.T, reader use_case.Reader, records []*use_case.RecordOrError) {
	for _, record := range records {
		switch record.InnerType {
		case use_case.RecordOrErrorInnerTypeError:
			if _, err := reader.Read(); err != record.Err {
				t.Errorf("error should happen here, expected: %+v, actual: %+v", record.Err, err)
			}
		case use_case.RecordOrErrorInnerTypeRecord:
			if _, err := reader.Read(); err != nil {
				t.Errorf("error should not happen here, expected: %+v, actual: %+v", nil, err)
			}
		}
	}

	if _, err := reader.Read(); err != io.EOF {
		t.Errorf("io.EOF should be returned here, expected: %+v, actual: %+v", io.EOF, err)
	}
}
