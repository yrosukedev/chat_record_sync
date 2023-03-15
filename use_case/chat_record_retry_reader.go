package use_case

import "github.com/yrosukedev/chat_record_sync/business"

type ChatRecordRetryReader struct {
	reader        ChatRecordReader
	maxRetryTimes uint
}

func NewChatRecordRetryReader(reader ChatRecordReader, maxRetryTimes uint) *ChatRecordRetryReader {
	return &ChatRecordRetryReader{
		reader:        reader,
		maxRetryTimes: maxRetryTimes,
	}
}

func (c *ChatRecordRetryReader) Read() (record *business.ChatRecord, err error) {
	return c.reader.Read()
}
