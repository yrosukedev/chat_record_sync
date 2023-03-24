package http_controller

type ResponseCode int

const (
	ResponseCodeOK                ResponseCode = 0
	ResponseCodeMarshalJsonFailed ResponseCode = 1
)

const (
	ResponseMsgOK                = "ok"
	ResponseMsgUnknownError      = "unknown error"
	ResponseMsgMarshalJsonFailed = "marshal json failed"
)

func (r ResponseCode) Msg() string {
	codeToMsgs := map[ResponseCode]string{
		ResponseCodeOK:                ResponseMsgOK,
		ResponseCodeMarshalJsonFailed: ResponseMsgMarshalJsonFailed,
	}

	if msg, ok := codeToMsgs[r]; ok {
		return msg
	}

	return ResponseMsgUnknownError
}
