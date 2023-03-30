package wecom_openapi_adapter

import (
	"context"
	"github.com/xen0n/go-workwx"
	wecom_chat2 "github.com/yrosukedev/chat_record_sync/chat_sync/wecom_chat"
	"github.com/yrosukedev/chat_record_sync/logger"
)

type WeComOpenAPIAdapter struct {
	ctx      context.Context
	wecomApp *workwx.WorkwxApp
	logger   logger.Logger
}

func NewWeComOpenAPIAdapter(ctx context.Context, wecomApp *workwx.WorkwxApp, logger logger.Logger) wecom_chat2.OpenAPIService {
	wecomApp.SpawnAccessTokenRefresher()
	return &WeComOpenAPIAdapter{
		ctx:      ctx,
		wecomApp: wecomApp,
		logger:   logger,
	}
}

func (w *WeComOpenAPIAdapter) GetUserInfoByID(id string) (userInfo *wecom_chat2.WeComUserInfo, err error) {
	w.logger.Info(w.ctx, "[wecom open api] will get user info, user id: %v", id)

	rawUserInfo, err := w.wecomApp.GetUser(id)
	if err != nil {
		w.logger.Error(w.ctx, "[wecom open api] fails to get user info, user id: %v, error: %v", id, err)
		return nil, err
	}

	userInfo = &wecom_chat2.WeComUserInfo{
		UserID: rawUserInfo.UserID,
		Name:   rawUserInfo.Name,
	}

	w.logger.Info(w.ctx, "[wecom open api] succeeds to get user info, user id: %v", id)

	return userInfo, nil
}

func (w *WeComOpenAPIAdapter) GetExternalContactByID(externalId string) (contact *wecom_chat2.WeComExternalContact, err error) {
	//TODO implement me
	panic("implement me")
}
