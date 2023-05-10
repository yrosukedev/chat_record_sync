package wecom

import (
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/chat_sync/reader/pagination"
	"math"
)

type PaginatedReaderAdapter struct {
	chatRecordService ChatRecordService
	recordTransformer RecordTransformer
}

func NewPaginatedReaderAdapter(
	chatRecordService ChatRecordService,
	recordTransformer RecordTransformer) *PaginatedReaderAdapter {
	return &PaginatedReaderAdapter{
		chatRecordService: chatRecordService,
		recordTransformer: recordTransformer,
	}
}

func (p *PaginatedReaderAdapter) Read(inPageToken *pagination.PageToken, pageSize uint64) (records []*business.ChatRecord, outPageToken *pagination.PageToken, err error) {
	outPageToken = inPageToken

	seq := p.seqFrom(inPageToken)

	wecomRecords, err := p.chatRecordService.Read(seq, pageSize)
	if err != nil {
		return nil, outPageToken, err
	}

	outPageToken = p.updatePageToken(inPageToken, wecomRecords)

	records, err = p.transformWecomRecords(wecomRecords)
	if err != nil {
		return nil, outPageToken, err
	}

	return records, outPageToken, nil
}

func (p *PaginatedReaderAdapter) seqFrom(inPageToken *pagination.PageToken) uint64 {
	seq := uint64(0)
	if inPageToken != nil {
		seq = inPageToken.Value
	}
	return seq
}

func (p *PaginatedReaderAdapter) transformWecomRecords(wecomRecords []*ChatRecord) (records []*business.ChatRecord, err error) {
	for _, wecomRecord := range wecomRecords {
		record, err := p.recordTransformer.Transform(wecomRecord)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}
	return records, nil
}

func (p *PaginatedReaderAdapter) updatePageToken(inPageToken *pagination.PageToken, wecomRecords []*ChatRecord) *pagination.PageToken {
	result := uint64(0)
	if inPageToken != nil {
		result = inPageToken.Value
	}

	for _, wecomRecord := range wecomRecords {
		result = uint64(math.Max(float64(result), float64(wecomRecord.Seq)))
	}

	if result > 0 {
		return pagination.NewPageToken(result)
	} else {
		return nil
	}
}
