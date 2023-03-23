package wecom_openapi_adapter

import (
	"github.com/xen0n/go-workwx"
	"github.com/yrosukedev/chat_record_sync/wecom_chat"
)

type WeComOpenAPIAdapter struct {
	wecomApp *workwx.WorkwxApp
}

func NewWeComOpenAPIAdapter(wecomApp *workwx.WorkwxApp) wecom_chat.OpenAPIService {
	wecomApp.SpawnAccessTokenRefresher()
	return &WeComOpenAPIAdapter{
		wecomApp: wecomApp,
	}
}

func (w *WeComOpenAPIAdapter) GetUserInfoByID(id string) (userInfo *wecom_chat.WeComUserInfo, err error) {
	rawUserInfo, err := w.wecomApp.GetUser(id)
	if err != nil {
		return nil, err
	}

	userInfo = &wecom_chat.WeComUserInfo{
		UserID: rawUserInfo.UserID,
		Name:   rawUserInfo.Name,
	}

	return userInfo, nil
}

func (w *WeComOpenAPIAdapter) GetExternalContactByID(externalId string) (contact *wecom_chat.WeComExternalContact, err error) {
	//TODO implement me
	panic("implement me")
}
