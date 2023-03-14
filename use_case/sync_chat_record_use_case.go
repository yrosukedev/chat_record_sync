package use_case

import (
	"context"
	"io"
)

type SyncChatRecordUseCase struct {
	reader ChatRecordReader
	writer ChatRecordWriter
}

func NewSyncChatRecordUseCase(reader ChatRecordReader, writer ChatRecordWriter) *SyncChatRecordUseCase {
	return &SyncChatRecordUseCase{
		reader: reader,
		writer: writer,
	}
}

func (u *SyncChatRecordUseCase) Run(ctx context.Context) []*SyncError {
	var errs []*SyncError

	for {
		record, err := u.reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			errs = append(errs, NewReaderError(err))
			continue
		}

		if err := u.writer.Write(record); err != nil {
			continue
		}
	}

	return errs
}
