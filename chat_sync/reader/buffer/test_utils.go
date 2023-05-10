package buffer

import (
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/chat_sync/use_case"
	"testing"
)

type recordsOrErr interface {
	records() []*business.ChatRecord
	err() error
}

type recordsWrapper struct {
	_records []*business.ChatRecord
}

func (r *recordsWrapper) records() []*business.ChatRecord {
	return r._records
}

func (r *recordsWrapper) err() error {
	return nil
}

type errWrapper struct {
	_err error
}

func (e *errWrapper) records() []*business.ChatRecord {
	return nil
}

func (e *errWrapper) err() error {
	return e._err
}

func givenRecordsToRefill(batchReader *MockBatchReader, recordsGroups [][]*business.ChatRecord) {
	groupIdx := 0
	batchReader.EXPECT().Read().DoAndReturn(func() ([]*business.ChatRecord, error) {
		defer func() { groupIdx += 1 }()
		return recordsGroups[groupIdx], nil
	}).Times(len(recordsGroups))
}

func givenRecordsOrErrorsToRefill(batchReader *MockBatchReader, recordsOrErrorGroups []recordsOrErr) {
	groupIdx := 0
	batchReader.EXPECT().Read().DoAndReturn(func() ([]*business.ChatRecord, error) {
		defer func() { groupIdx += 1 }()
		switch recordsOrError := recordsOrErrorGroups[groupIdx].(type) {
		case *recordsWrapper:
			return recordsOrError.records(), nil
		case *errWrapper:
			return nil, recordsOrError.err()
		default:
			return nil, nil
		}
	}).Times(len(recordsOrErrorGroups))
}

func expectReaderToReadRecords(t *testing.T, reader *Reader, records []*business.ChatRecord) {
	for _, expected := range records {
		actual, err := reader.Read()
		if err != nil {
			t.Errorf("error should not happen here, expected: %+v, actual: %+v", nil, err)
			return
		}

		if expected != actual {
			t.Errorf("records are not matched, expected: %+v, actual: %+v", expected, actual)
			return
		}
	}
}

func expectReaderToReadRecordOrError(t *testing.T, reader *Reader, recordsOrErrors []use_case.RecordOrError) {
	for _, expected := range recordsOrErrors {
		actual, err := reader.Read()
		switch expected := expected.(type) {
		case use_case.RecordWrapper:
			if expected.Record() != actual {
				t.Errorf("records are not matched, expected: %+v, actual: %+v", expected.Record(), actual)
				return
			}
		case use_case.ErrorWrapper:
			if expected.Error() != err {
				t.Errorf("errors are not matched, expected: %+v, actual: %+v", expected.Error(), err)
				return
			}
		}
	}
}
