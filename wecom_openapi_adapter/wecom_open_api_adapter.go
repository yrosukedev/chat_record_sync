package wecom_openapi_adapter

import (
	"context"
	"github.com/xen0n/go-workwx"
	"github.com/yrosukedev/chat_record_sync/logger"
	"github.com/yrosukedev/chat_record_sync/wecom_chat"
)

type WeComOpenAPIAdapter struct {
	ctx      context.Context
	wecomApp *workwx.WorkwxApp
	logger   logger.Logger
}

func NewWeComOpenAPIAdapter(ctx context.Context, wecomApp *workwx.WorkwxApp, logger logger.Logger) wecom_chat.OpenAPIService {
	wecomApp.SpawnAccessTokenRefresher()
	return &WeComOpenAPIAdapter{
		ctx:      ctx,
		wecomApp: wecomApp,
		logger:   logger,
	}
}

func (w *WeComOpenAPIAdapter) GetUserInfoByID(id string) (userInfo *wecom_chat.WeComUserInfo, err error) {
	w.logger.Info(w.ctx, "[wecom open api] will get user info, user id: %v", id)

	rawUserInfo, err := w.wecomApp.GetUser(id)
	if err != nil {
		w.logger.Error(w.ctx, "[wecom open api] fails to get user info, user id: %v, error: %v", id, err)
		return nil, err
	}

	userInfo = &wecom_chat.WeComUserInfo{
		UserID: rawUserInfo.UserID,
		Name:   rawUserInfo.Name,
	}

	w.logger.Info(w.ctx, "[wecom open api] succeeds to get user info, user id: %v", id)

	return userInfo, nil
}

func (w *WeComOpenAPIAdapter) GetExternalContactByID(externalId string) (contact *wecom_chat.WeComExternalContact, err error) {
	//TODO implement me
	panic("implement me")
}
