package wecom_chat

import (
	"github.com/yrosukedev/chat_record_sync/business"
	"github.com/yrosukedev/chat_record_sync/paginated_reader"
	"math"
)

type PaginatedBufferedReaderAdapter struct {
	chatRecordService ChatRecordService
	openAPIService    OpenAPIService
	transformer       ChatRecordTransformer
}

func NewPaginatedBufferedReaderAdapter(
	chatRecordService ChatRecordService,
	openAPIService OpenAPIService,
	transformer ChatRecordTransformer) *PaginatedBufferedReaderAdapter {
	return &PaginatedBufferedReaderAdapter{
		chatRecordService: chatRecordService,
		openAPIService:    openAPIService,
		transformer:       transformer,
	}
}

func (p *PaginatedBufferedReaderAdapter) Read(inPageToken *paginated_reader.PageToken, pageSize uint64) (records []*business.ChatRecord, outPageToken *paginated_reader.PageToken, err error) {
	outPageToken = inPageToken

	wecomRecords, err := p.chatRecordService.Read(inPageToken.Value, pageSize)
	if err != nil {
		return nil, outPageToken, err
	}

	for _, wecomRecord := range wecomRecords {
		user, err := p.openAPIService.GetUserInfoByID(wecomRecord.From)
		if err != nil {
			return nil, outPageToken, err
		}

		var contacts []*WeComExternalContact
		for _, contactId := range wecomRecord.ToList {
			contact, err := p.openAPIService.GetExternalContactByID(contactId)
			if err != nil {
				return nil, outPageToken, err
			}
			contacts = append(contacts, contact)
		}

		record, err := p.transformer.Transform(wecomRecord, user, contacts)
		if err != nil {
			return nil, outPageToken, err
		}

		records = append(records, record)

		outPageToken = paginated_reader.NewPageToken(uint64(math.Max(float64(inPageToken.Value), float64(wecomRecord.Seq))))
	}

	return records, outPageToken, nil
}
