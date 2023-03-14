package use_case

import (
	"encoding/json"
	"fmt"
	"github.com/yrosukedev/chat_record_sync/business"
)

type SyncErrorType = string

var (
	SyncErrorTypeRead  SyncErrorType = "Read"
	SyncErrorTypeWrite SyncErrorType = "Write"
)

type SyncError struct {
	ErrorType SyncErrorType        // The type of the error
	Err       error                // The actual error
	Record    *business.ChatRecord // The chat record. If the error type is SyncErrorTypeRead, this field is nil
}

func NewReaderError(err error) *SyncError {
	return &SyncError{
		ErrorType: SyncErrorTypeRead,
		Err:       err,
	}
}

func NewWriterError(err error, record *business.ChatRecord) *SyncError {
	return &SyncError{
		ErrorType: SyncErrorTypeWrite,
		Err:       err,
		Record:    record,
	}
}

func (s *SyncError) Error() string {
	marshaledRecord := ""
	if s.Record != nil {
		bs, err := json.Marshal(s.Record)
		if err != nil {
			marshaledRecord = fmt.Sprintf("fail to marshal record(error: %v)", err)
		} else {
			marshaledRecord = string(bs)
		}
	}

	return fmt.Sprintf("fail to sync chat record, error type: %v, actual error: %v, record: '%v'", s.ErrorType, s.Err, marshaledRecord)
}
