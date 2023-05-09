package chat_record

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"testing"
)

func TestFieldsWriter_Write_ZeroField(t *testing.T) {
	// if the fields is empty, the storage shouldn't be called.

	// Given
	ctrl := gomock.NewController(t)
	fieldsFormatter := NewMockFieldsFormatter(ctrl)
	fieldsStorage := NewMockFieldsStorage(ctrl)
	fieldsWriter := NewFieldsWriter(fieldsFormatter, fieldsStorage)

	chatRecord := &business.ChatRecord{}
	fields := map[string]interface{}{}

	fieldsFormatter.EXPECT().Format(gomock.Eq(chatRecord)).Return(fields, nil).Times(1)
	fieldsStorage.EXPECT().Write(gomock.Any(), gomock.Any()).Times(0)

	// When
	err := fieldsWriter.Write(chatRecord, "26df9a5c-55d8-4c52-b6ce-203325568178")

	// Then
	assert.NoError(t, err)
}

func TestFieldsWriter_Write_NilField(t *testing.T) {
	// if the fields is nil, the storage shouldn't be called.

	// Given
	ctrl := gomock.NewController(t)
	fieldsFormatter := NewMockFieldsFormatter(ctrl)
	fieldsStorage := NewMockFieldsStorage(ctrl)
	fieldsWriter := NewFieldsWriter(fieldsFormatter, fieldsStorage)

	chatRecord := &business.ChatRecord{}

	fieldsFormatter.EXPECT().Format(gomock.Eq(chatRecord)).Return(nil, nil).Times(1)
	fieldsStorage.EXPECT().Write(gomock.Any(), gomock.Any()).Times(0)

	// When
	err := fieldsWriter.Write(chatRecord, "26df9a5c-55d8-4c52-b6ce-203325568178")

	// Then
	assert.NoError(t, err)
}

func TestFieldsWriter_Write_NilChatRecord(t *testing.T) {
	// if the chatRecord is nil, neither the formatter nor the storage should be called,
	// and the writer should return an error.

	// Given
	ctrl := gomock.NewController(t)
	fieldsFormatter := NewMockFieldsFormatter(ctrl)
	fieldsStorage := NewMockFieldsStorage(ctrl)
	fieldsWriter := NewFieldsWriter(fieldsFormatter, fieldsStorage)

	fieldsFormatter.EXPECT().Format(gomock.Any()).Times(0)
	fieldsStorage.EXPECT().Write(gomock.Any(), gomock.Any()).Times(0)

	// When
	err := fieldsWriter.Write(nil, "26df9a5c-55d8-4c52-b6ce-203325568178")

	// Then
	assert.Error(t, err)
}

func TestFieldsWriter_Write_OneField(t *testing.T) {
	// if the size of the fields is 1, the storage should be called once.

	// Given
	ctrl := gomock.NewController(t)
	fieldsFormatter := NewMockFieldsFormatter(ctrl)
	fieldsStorage := NewMockFieldsStorage(ctrl)
	fieldsWriter := NewFieldsWriter(fieldsFormatter, fieldsStorage)

	chatRecord := &business.ChatRecord{}
	fields := map[string]interface{}{
		"field1": "value1",
	}

	fieldsFormatter.EXPECT().Format(gomock.Eq(chatRecord)).Return(fields, nil).Times(1)
	fieldsStorage.EXPECT().Write(gomock.Any(), gomock.Any()).Times(1)

	// When
	err := fieldsWriter.Write(chatRecord, "26df9a5c-55d8-4c52-b6ce-203325568178")

	// Then
	assert.NoError(t, err)
}

func TestFieldsWriter_Write_MultipleFields(t *testing.T) {
	// if the size of the fields is greater than 1, the storage should be called once.

	// Given
	ctrl := gomock.NewController(t)
	fieldsFormatter := NewMockFieldsFormatter(ctrl)
	fieldsStorage := NewMockFieldsStorage(ctrl)
	fieldsWriter := NewFieldsWriter(fieldsFormatter, fieldsStorage)

	chatRecord := &business.ChatRecord{}
	fields := map[string]interface{}{
		"field1": "value1",
		"field2": "value2",
		"field3": "value3",
	}

	fieldsFormatter.EXPECT().Format(gomock.Eq(chatRecord)).Return(fields, nil).Times(1)
	fieldsStorage.EXPECT().Write(gomock.Any(), gomock.Any()).Times(1)

	// When
	err := fieldsWriter.Write(chatRecord, "26df9a5c-55d8-4c52-b6ce-203325568178")

	// Then
	assert.NoError(t, err)
}
