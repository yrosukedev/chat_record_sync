package buffer_reader

import (
	"github.com/yrosukedev/chat_record_sync/business"
	"io"
)

type ChatRecordBufferedReaderAdapter struct {
	bufferedReader  ChatRecordBufferedReader
	idx             int
	bufferedRecords []*business.ChatRecord
}

func NewChatRecordBufferedReaderAdapter(bufferedReader ChatRecordBufferedReader) *ChatRecordBufferedReaderAdapter {
	return &ChatRecordBufferedReaderAdapter{
		bufferedReader: bufferedReader,
		idx:            0,
	}
}

func (c *ChatRecordBufferedReaderAdapter) Read() (record *business.ChatRecord, err error) {
	if err := c.refillBufferIfNeeded(); err != nil {
		return nil, err
	}

	return c.readFromBuffer()
}

func (c *ChatRecordBufferedReaderAdapter) refillBufferIfNeeded() error {
	if !c.needRefill() {
		return nil
	}

	var err error
	c.bufferedRecords, err = c.bufferedReader.Read()
	c.idx = 0
	return err
}

func (c *ChatRecordBufferedReaderAdapter) needRefill() bool {
	return c.idx >= len(c.bufferedRecords)
}

func (c *ChatRecordBufferedReaderAdapter) readFromBuffer() (*business.ChatRecord, error) {
	if c.needRefill() { // still need refill, but already refilled, that means reaching the end
		return nil, io.EOF
	}

	defer func() { c.idx += 1 }()
	return c.bufferedRecords[c.idx], nil
}
