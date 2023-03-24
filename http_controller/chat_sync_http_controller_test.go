package http_controller

import (
	"context"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"net/http"
	"strings"
	"testing"
)

func TestSucceeds(t *testing.T) {
	// Given
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	useCase := NewMockUseCase(ctrl)
	controller := NewChatSyncHTTPController(ctx, useCase)
	responseWriter := NewMockResponseWriter(ctrl)
	chatSyncReponse := ChatSyncResponse{
		Code: ResponseCodeOK,
		Msg:  ResponseMsgOK,
	}

	chatSyncReponseJson, err := json.Marshal(chatSyncReponse)
	if err != nil {
		t.Errorf("error shouldn't happen here, expected: %v, actual: %v", nil, err)
		return
	}

	// Then
	useCase.EXPECT().Run(gomock.Eq(ctx)).Return(nil).Times(1)

	responseWriter.EXPECT().WriteHeader(gomock.Eq(http.StatusOK)).Times(1)
	responseWriter.EXPECT().Write(gomock.Eq(chatSyncReponseJson)).Return(len(chatSyncReponseJson), nil).Times(1)

	// When
	request, err := http.NewRequest(http.MethodPost, "/chat_sync", strings.NewReader(""))
	if err != nil {
		t.Errorf("error shouldn't happen here, expected: %v, actual: %v", nil, err)
		return
	}

	controller.ServeHTTP(responseWriter, request)
}
