package wecom_chat

type OpenAPIService interface {
	GetUserInfoByID(id string) (userInfo *WeComUserInfo, err error)
	GetExternalContactByID(externalId string) (contact *WeComExternalContact, err error)
}
