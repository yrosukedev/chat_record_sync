package transformer

type UserNameFetcher struct {
	openAPIService OpenAPIService
}

func NewUserNameFetcher(openAPIService OpenAPIService) *UserNameFetcher {
	return &UserNameFetcher{
		openAPIService: openAPIService,
	}
}

func (f *UserNameFetcher) FetchName(id string) (name string, err error) {
	userInfo, err := f.openAPIService.GetUserInfoByID(id)

	if err != nil {
		return "", err
	}

	return userInfo.Name, nil
}
