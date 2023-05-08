package transformer

import "github.com/yrosukedev/chat_record_sync/chat_sync/wecom"

type OpenAPIService interface {
	GetUserInfoByID(id string) (userInfo *wecom.UserInfo, err error)
	GetExternalContactByID(externalId string) (contact *wecom.ExternalContact, err error)
}

type MsgAuditOpenAPIService interface {
	GetInternalRoomByID(roomId string) (room *wecom.InternalRoom, err error)
}
