package transformer

type ContactNameFetcher struct {
	openAPIService OpenAPIService
}

func NewContactNameFetcher(openAPIService OpenAPIService) *ContactNameFetcher {
	return &ContactNameFetcher{
		openAPIService: openAPIService,
	}
}

func (f *ContactNameFetcher) FetchName(id string) (name string, err error) {
	userInfo, err := f.openAPIService.GetExternalContactByID(id)

	if err != nil {
		return "", err
	}

	return userInfo.Name, nil
}
