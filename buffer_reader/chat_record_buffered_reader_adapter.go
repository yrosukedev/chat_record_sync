package buffer_reader

import (
	"github.com/yrosukedev/chat_record_sync/business"
)

type ChatRecordBufferedReaderAdapter struct {
	bufferedReader ChatRecordBufferedReader
}

func NewChatRecordBufferedReaderAdapter(bufferedReader ChatRecordBufferedReader) *ChatRecordBufferedReaderAdapter {
	return &ChatRecordBufferedReaderAdapter{
		bufferedReader: bufferedReader,
	}
}

func (c *ChatRecordBufferedReaderAdapter) Read() (record *business.ChatRecord, err error) {
	records, err := c.bufferedReader.Read()
	return records[0], err
}
