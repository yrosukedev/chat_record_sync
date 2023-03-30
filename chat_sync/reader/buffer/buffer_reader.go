package buffer

import (
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"io"
)

type Reader struct {
	batchReader     BatchReader
	idx             int
	bufferedRecords []*business.ChatRecord
}

func NewReader(batchReader BatchReader) *Reader {
	return &Reader{
		batchReader: batchReader,
		idx:         0,
	}
}

func (c *Reader) Read() (record *business.ChatRecord, err error) {
	if err := c.refillBufferIfNeeded(); err != nil {
		return nil, err
	}

	return c.readFromBuffer()
}

func (c *Reader) refillBufferIfNeeded() error {
	if !c.needRefill() {
		return nil
	}

	var err error
	c.bufferedRecords, err = c.batchReader.Read()
	c.idx = 0
	return err
}

func (c *Reader) needRefill() bool {
	return c.idx >= len(c.bufferedRecords)
}

func (c *Reader) readFromBuffer() (*business.ChatRecord, error) {
	if c.needRefill() { // still need refill, but already refilled, that means reaching the end
		return nil, io.EOF
	}

	defer func() { c.idx += 1 }()
	return c.bufferedRecords[c.idx], nil
}
