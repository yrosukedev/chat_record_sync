package openapi

import (
	"context"
	"github.com/xen0n/go-workwx"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom"
	"github.com/yrosukedev/chat_record_sync/logger"
)

type MsgAuditOpenAPIAdapter struct {
	ctx      context.Context
	wecomApp *workwx.WorkwxApp
	logger   logger.Logger
}

func NewMsgAuditOpenAPIAdapter(ctx context.Context, wecomApp *workwx.WorkwxApp, logger logger.Logger) *MsgAuditOpenAPIAdapter {
	return &MsgAuditOpenAPIAdapter{
		ctx:      ctx,
		wecomApp: wecomApp,
		logger:   logger,
	}
}

func (w *MsgAuditOpenAPIAdapter) GetInternalRoomByID(roomId string) (room *wecom.InternalRoom, err error) {
	w.logger.Info(w.ctx, "[wecom open api] will get internal room info, room id: %v", roomId)

	groupChat, err := w.wecomApp.GetMsgAuditGroupChat(roomId)
	if err != nil {
		w.logger.Error(w.ctx, "[wecom open api] fails to get internal room info, room id: %v, error: %v", roomId, err)
		return nil, err
	}

	room = &wecom.InternalRoom{
		RoomID: roomId,
		Name:   groupChat.RoomName,
	}

	w.logger.Info(w.ctx, "[wecom open api] succeeds to get internal room info, room id: %v", roomId)

	return room, nil
}
