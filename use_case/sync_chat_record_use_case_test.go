package use_case

import (
	"context"
	"github.com/golang/mock/gomock"
	"io"
	"testing"
)

func TestZeroRecord(t *testing.T) {
	// Given
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	reader := NewMockChatRecordReader(ctrl)
	writer := NewMockChatRecordWriter(ctrl)
	useCase := NewSyncChatRecordUseCase(reader, writer)

	reader.EXPECT().Read().Times(1).Return(nil, io.EOF)

	// Then
	writer.EXPECT().Write(gomock.Any()).Times(0)

	// When
	useCase.Run(ctx)
}
