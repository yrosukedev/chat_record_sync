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

func (u *SyncChatRecordUseCase) Run(ctx context.Context) error {

	for {
		record, err := u.reader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			continue
		}

		if err := u.writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}
