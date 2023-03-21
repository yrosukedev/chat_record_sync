package wecom_chat

import (
	"github.com/yrosukedev/chat_record_sync/business"
	"github.com/yrosukedev/chat_record_sync/paginated_reader"
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
	_, err = p.chatRecordService.Read(inPageToken.Value, pageSize)
	if err != nil {
		return nil, inPageToken, err
	}

	return nil, inPageToken, nil
}
