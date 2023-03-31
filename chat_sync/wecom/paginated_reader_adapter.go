package wecom

import (
	"github.com/yrosukedev/chat_record_sync/chat_sync/business"
	"github.com/yrosukedev/chat_record_sync/chat_sync/reader/pagination"
	"math"
)

type PaginatedReaderAdapter struct {
	chatRecordService ChatRecordService
	openAPIService    OpenAPIService
	transformer       ChatRecordTransformer
}

func NewPaginatedReaderAdapter(
	chatRecordService ChatRecordService,
	openAPIService OpenAPIService,
	transformer ChatRecordTransformer) *PaginatedReaderAdapter {
	return &PaginatedReaderAdapter{
		chatRecordService: chatRecordService,
		openAPIService:    openAPIService,
		transformer:       transformer,
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

func (p *PaginatedReaderAdapter) getUserInfo(wecomRecord *ChatRecord) (*UserInfo, error) {
	if p.openAPIService == nil {
		return nil, nil
	}

	user, err := p.openAPIService.GetUserInfoByID(wecomRecord.From)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (p *PaginatedReaderAdapter) getContacts(contactIds []string) ([]*ExternalContact, error) {
	if p.openAPIService == nil {
		return nil, nil
	}

	var contacts []*ExternalContact
	for _, contactId := range contactIds {
		contact, err := p.openAPIService.GetExternalContactByID(contactId)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}
	return contacts, nil
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
