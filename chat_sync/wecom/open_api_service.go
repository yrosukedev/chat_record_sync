package wecom

type OpenAPIService interface {
	GetUserInfoByID(id string) (userInfo *UserInfo, err error)
	GetExternalContactByID(externalId string) (contact *ExternalContact, err error)
}
