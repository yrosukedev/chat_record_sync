package use_case

import (
	"github.com/yrosukedev/chat_record_sync/business"
	"io"
)

type ChatRecordRetryReader struct {
	reader                  ChatRecordReader
	maxRetryTimes           uint
	consecutiveFailureTimes uint
}

func NewChatRecordRetryReader(reader ChatRecordReader, maxRetryTimes uint) *ChatRecordRetryReader {
	return &ChatRecordRetryReader{
		reader:                  reader,
		maxRetryTimes:           maxRetryTimes,
		consecutiveFailureTimes: 0,
	}
}

// Read from the proxy reader and forward the result.
// If the number of proxy reader consecutive failure exceeds maxRetryTimes,
// stop forwarding the result and append io.EOF to indicate the end.
func (c *ChatRecordRetryReader) Read() (record *business.ChatRecord, err error) {
	if c.consecutiveFailureTimes > c.maxRetryTimes {
		return nil, io.EOF
	}

	record, err = c.reader.Read()

	if err != nil && err != io.EOF {
		c.consecutiveFailureTimes += 1
	}

	return record, err
}
