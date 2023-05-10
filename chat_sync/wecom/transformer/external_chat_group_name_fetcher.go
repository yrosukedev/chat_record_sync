package transformer

type ExternalChatGroupNameFetcher struct {
	openAPIService OpenAPIService
}

func NewExternalChatGroupNameFetcher(openAPIService OpenAPIService) *ExternalChatGroupNameFetcher {
	return &ExternalChatGroupNameFetcher{
		openAPIService: openAPIService,
	}
}

func (f *ExternalChatGroupNameFetcher) FetchName(id string) (name string, err error) {
	room, err := f.openAPIService.GetExternalRoomByID(id)

	if err != nil {
		return "", err
	}

	return room.Name, nil
}
