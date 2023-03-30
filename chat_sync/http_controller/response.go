package http_controller

type ChatSyncResponse struct {
	Code ResponseCode `json:"code"`
	Msg  string       `json:"msg"`
}
