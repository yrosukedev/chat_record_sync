package retry

import (
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/chat_sync/use_case"
	"io"
	"math"
)

type Reader struct {
	reader                  use_case.Reader
	maxRetryTimes           uint
	consecutiveFailureTimes uint
}

// NewRetryReader act as an adapter for the original reader.
//
// maxRetryTimes should be greater than or equal to 1. if it is less than 1, 1 will be used.
func NewRetryReader(reader use_case.Reader, maxRetryTimes uint) *Reader {
	return &Reader{
		reader:                  reader,
		maxRetryTimes:           uint(math.Max(float64(maxRetryTimes), 1)),
		consecutiveFailureTimes: 0,
	}
}

// Read from the proxy reader and forward the result.
// If the number of proxy reader consecutive failure exceeds maxRetryTimes,
// stop forwarding the result and append io.EOF to indicate the end.
func (c *Reader) Read() (record *business.ChatRecord, err error) {
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
