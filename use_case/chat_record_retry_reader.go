package use_case

import (
	"github.com/yrosukedev/chat_record_sync/business"
	"io"
	"math"
)

type ChatRecordRetryReader struct {
	reader                  ChatRecordReader
	maxRetryTimes           uint
	consecutiveFailureTimes uint
}

// NewChatRecordRetryReader act as an adapter for the original reader.
//
// maxRetryTimes should be greater than or equal to 1. if it is less than 1, 1 will be used.
func NewChatRecordRetryReader(reader ChatRecordReader, maxRetryTimes uint) *ChatRecordRetryReader {
	return &ChatRecordRetryReader{
		reader:                  reader,
		maxRetryTimes:           uint(math.Max(float64(maxRetryTimes), 1)),
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

	if err == nil || err == io.EOF {
		c.consecutiveFailureTimes = 0
	} else {
		c.consecutiveFailureTimes += 1
	}

	if c.consecutiveFailureTimes > c.maxRetryTimes {
		return nil, io.EOF
	}

	return record, err
}
