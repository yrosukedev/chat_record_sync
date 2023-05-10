package transformer

type InternalChatGroupNameFetcher struct {
	openAPIService MsgAuditOpenAPIService
}

func NewInternalChatGroupNameFetcher(openAPIService MsgAuditOpenAPIService) *InternalChatGroupNameFetcher {
	return &InternalChatGroupNameFetcher{
		openAPIService: openAPIService,
	}
}

func (f *InternalChatGroupNameFetcher) FetchName(id string) (name string, err error) {
	room, err := f.openAPIService.GetInternalRoomByID(id)

	if err != nil {
		return "", err
	}

	return room.Name, nil
}
