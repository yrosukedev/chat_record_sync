package openapi

import (
	"context"
	"github.com/xen0n/go-workwx"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom/transformer"
	"github.com/yrosukedev/chat_record_sync/logger"
)

type Adapter struct {
	ctx      context.Context
	wecomApp *workwx.WorkwxApp
	logger   logger.Logger
}

func NewAdapter(ctx context.Context, wecomApp *workwx.WorkwxApp, logger logger.Logger) transformer.OpenAPIService {
	wecomApp.SpawnAccessTokenRefresher()
	return &Adapter{
		ctx:      ctx,
		wecomApp: wecomApp,
		logger:   logger,
	}
}

// GetUserInfoByID get user info by user id
// @param id user id
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
	w.logger.Info(w.ctx, "[wecom open api] will get external contact info, external id: %v", externalId)

	rawContact, err := w.wecomApp.GetExternalContact(externalId)
	if err != nil {
		w.logger.Error(w.ctx, "[wecom open api] fails to get external contact info, external id: %v, error: %v", externalId, err)
		return nil, err
	}

	contact = &wecom.ExternalContact{
		ExternalUserID: rawContact.ExternalContact.ExternalUserid,
		Name:           rawContact.ExternalContact.Name,
	}

	w.logger.Info(w.ctx, "[wecom open api] succeeds to get external contact info, external id: %v", externalId)

	return contact, nil
}

func (w *Adapter) GetExternalRoomByID(roomId string) (room *wecom.ExternalRoom, err error) {
	w.logger.Info(w.ctx, "[wecom open api] will get external room info, room id: %v", roomId)

	chatInfo, err := w.wecomApp.GetAppChatInfo(roomId)
	if err != nil {
		w.logger.Error(w.ctx, "[wecom open api] fails to get external room info, room id: %v, error: %v", roomId, err)
		return nil, err
	}

	room = &wecom.ExternalRoom{
		RoomID: roomId,
		Name:   chatInfo.Name,
	}

	w.logger.Info(w.ctx, "[wecom open api] succeeds to get external room info, room id: %v", roomId)

	return room, nil
}
