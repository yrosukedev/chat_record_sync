package wecom_chat

import (
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/chat_sync/reader/pagination"
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

func (p *PaginatedBufferedReaderAdapter) Read(inPageToken *pagination.PageToken, pageSize uint64) (records []*business.ChatRecord, outPageToken *pagination.PageToken, err error) {
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

func (p *PaginatedBufferedReaderAdapter) seqFrom(inPageToken *pagination.PageToken) uint64 {
	seq := uint64(0)
	if inPageToken != nil {
		seq = inPageToken.Value
	}
	return seq
}

func (p *PaginatedBufferedReaderAdapter) transformWecomRecords(wecomRecords []*WeComChatRecord) (records []*business.ChatRecord, err error) {
	for _, wecomRecord := range wecomRecords {
		user, err := p.getUserInfo(wecomRecord)
		if err != nil {
			return nil, err
		}

		contacts, err := p.getContacts(wecomRecord.ToList)
		if err != nil {
			return nil, err
		}

		record, err := p.transformer.Transform(wecomRecord, user, contacts)
		if err != nil {
			return nil, err
		}

		records = append(records, record)
	}
	return records, nil
}

func (p *PaginatedBufferedReaderAdapter) getUserInfo(wecomRecord *WeComChatRecord) (*WeComUserInfo, error) {
	if p.openAPIService == nil {
		return nil, nil
	}

	user, err := p.openAPIService.GetUserInfoByID(wecomRecord.From)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (p *PaginatedBufferedReaderAdapter) getContacts(contactIds []string) ([]*WeComExternalContact, error) {
	if p.openAPIService == nil {
		return nil, nil
	}

	var contacts []*WeComExternalContact
	for _, contactId := range contactIds {
		contact, err := p.openAPIService.GetExternalContactByID(contactId)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}
	return contacts, nil
}

func (p *PaginatedBufferedReaderAdapter) updatePageToken(inPageToken *pagination.PageToken, wecomRecords []*WeComChatRecord) *pagination.PageToken {
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
