package use_case

import (
	"context"
	"io"
)

type ChatSyncUseCase struct {
	reader Reader
	writer Writer
}

func NewChatSyncUseCase(reader Reader, writer Writer) *ChatSyncUseCase {
	return &ChatSyncUseCase{
		reader: reader,
		writer: writer,
	}
}

func (u *ChatSyncUseCase) Run(ctx context.Context) []*SyncError {
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
			errs = append(errs, NewWriterError(err, record))
			continue
		}
	}

	return errs
}
