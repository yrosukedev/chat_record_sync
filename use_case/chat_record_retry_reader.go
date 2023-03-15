package use_case

import (
	"github.com/yrosukedev/chat_record_sync/business"
	"io"
)

type ChatRecordRetryReader struct {
	reader        ChatRecordReader
	maxRetryTimes uint
	retriedTimes  uint
}

func NewChatRecordRetryReader(reader ChatRecordReader, maxRetryTimes uint) *ChatRecordRetryReader {
	return &ChatRecordRetryReader{
		reader:        reader,
		maxRetryTimes: maxRetryTimes,
		retriedTimes:  0,
	}
}

func (c *ChatRecordRetryReader) Read() (record *business.ChatRecord, err error) {
	for c.retriedTimes < c.maxRetryTimes {
		record, err := c.reader.Read()

		if err == nil || err == io.EOF {
			return record, err
		}

		c.retriedTimes += 1

		// TODO: log the error
	}

	return nil, io.EOF
}
