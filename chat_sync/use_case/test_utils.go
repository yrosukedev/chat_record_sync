package use_case

import (
	"github.com/golang/mock/gomock"
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"io"
)

type RecordOrError interface {
	Record() *business.ChatRecord
	Error() error
}

type RecordWrapper struct {
	record *business.ChatRecord
}

func (r RecordWrapper) Record() *business.ChatRecord {
	return r.record
}

func (r RecordWrapper) Error() error {
	return nil
}

type ErrorWrapper struct {
	err error
}

func (e ErrorWrapper) Record() *business.ChatRecord {
	return nil
}

func (e ErrorWrapper) Error() error {
	return e.err
}

func NewRecordOrErrorWithRecord(record *business.ChatRecord) RecordOrError {
	return RecordWrapper{
		record: record,
	}
}

func NewRecordOrErrorWithError(err error) RecordOrError {
	return ErrorWrapper{
		err: err,
	}
}

func ExpectRecordsToWrite(writer *MockWriter, records []*business.ChatRecord) {
	for _, r := range records {
		writer.
			EXPECT().
			Write(gomock.Eq(r)).
			Times(1)
	}
}

func GivenRecordsToRead(reader *MockReader, records []*business.ChatRecord) {
	idx := 0
	reader.
		EXPECT().
		Read().
		DoAndReturn(func() (*business.ChatRecord, error) {
			if idx >= len(records) {
				return nil, io.EOF
			}
			defer func() { idx += 1 }()
			return records[idx], nil
		}).
		AnyTimes()
}

func EncounterErrorWhileReadingRecords(reader *MockReader, records []RecordOrError) {
	idx := 0
	reader.
		EXPECT().
		Read().
		DoAndReturn(func() (*business.ChatRecord, error) {
			if idx >= len(records) {
				return nil, io.EOF
			}

			defer func() { idx += 1 }()

			switch record := records[idx].(type) {
			case ErrorWrapper:
				return nil, record.Error()
			case RecordWrapper:
				return record.Record(), nil
			default:
				return nil, nil
			}
		}).
		AnyTimes()
}
