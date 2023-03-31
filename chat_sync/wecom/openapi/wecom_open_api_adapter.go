package openapi

import (
	"context"
	"github.com/xen0n/go-workwx"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom"
	"github.com/yrosukedev/chat_record_sync/logger"
)

type Adapter struct {
	ctx      context.Context
	wecomApp *workwx.WorkwxApp
	logger   logger.Logger
}

func NewAdapter(ctx context.Context, wecomApp *workwx.WorkwxApp, logger logger.Logger) wecom.OpenAPIService {
	wecomApp.SpawnAccessTokenRefresher()
	return &Adapter{
		ctx:      ctx,
		wecomApp: wecomApp,
		logger:   logger,
	}
}

func (w *Adapter) GetUserInfoByID(id string) (userInfo *wecom.UserInfo, err error) {
	w.logger.Info(w.ctx, "[wecom open api] will get user info, user id: %v", id)

	rawUserInfo, err := w.wecomApp.GetUser(id)
	if err != nil {
		w.logger.Error(w.ctx, "[wecom open api] fails to get user info, user id: %v, error: %v", id, err)
		return nil, err
	}

	userInfo = &wecom.UserInfo{
		UserID: rawUserInfo.UserID,
		Name:   rawUserInfo.Name,
	}

	w.logger.Info(w.ctx, "[wecom open api] succeeds to get user info, user id: %v", id)

	return userInfo, nil
}

func (w *Adapter) GetExternalContactByID(externalId string) (contact *wecom.ExternalContact, err error) {
	//TODO implement me
	panic("implement me")
}
