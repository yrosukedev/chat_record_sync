package use_case

import (
	"github.com/golang/mock/gomock"
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"io"
)

type RecordOrErrorInnerType = string

var (
	RecordOrErrorInnerTypeRecord RecordOrErrorInnerType = "Record"
	RecordOrErrorInnerTypeError  RecordOrErrorInnerType = "Error"
)

type RecordOrError struct {
	InnerType RecordOrErrorInnerType
	Record    *business.ChatRecord
	Err       error
}

func NewRecordOrErrorWithRecord(record *business.ChatRecord) *RecordOrError {
	return &RecordOrError{
		InnerType: RecordOrErrorInnerTypeRecord,
		Record:    record,
	}
}

func NewRecordOrErrorWithError(err error) *RecordOrError {
	return &RecordOrError{
		InnerType: RecordOrErrorInnerTypeError,
		Err:       err,
	}
}

func ExpectRecordsToWrite(writer *MockChatRecordWriter, records []*business.ChatRecord) {
	for _, r := range records {
		writer.
			EXPECT().
			Write(gomock.Eq(r)).
			Times(1)
	}
}

func GivenRecordsToRead(reader *MockChatRecordReader, records []*business.ChatRecord) {
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

func EncounterErrorWhileReadingRecords(reader *MockChatRecordReader, records []*RecordOrError) {
	idx := 0
	reader.
		EXPECT().
		Read().
		DoAndReturn(func() (*business.ChatRecord, error) {
			if idx >= len(records) {
				return nil, io.EOF
			}

			defer func() { idx += 1 }()

			record := records[idx]
			switch record.InnerType {
			case RecordOrErrorInnerTypeError:
				return nil, record.Err
			default:
				return record.Record, nil
			}
		}).
		AnyTimes()
}
