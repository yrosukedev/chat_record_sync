package http_controller

type ResponseCode int

const (
	ResponseCodeOK                ResponseCode = 0
	ResponseCodeMarshalJsonFailed ResponseCode = 1
	ResponseCodeFailure           ResponseCode = 2
)

const (
	ResponseMsgOK                = "ok"
	ResponseMsgUnknownError      = "unknown error"
	ResponseMsgMarshalJsonFailed = "marshal json failed"
	ResponseMsgFailure           = "failed"
)

func (r ResponseCode) Msg() string {
	codeToMsgs := map[ResponseCode]string{
		ResponseCodeOK:                ResponseMsgOK,
		ResponseCodeMarshalJsonFailed: ResponseMsgMarshalJsonFailed,
		ResponseCodeFailure:           ResponseMsgFailure,
	}

	if msg, ok := codeToMsgs[r]; ok {
		return msg
	}

	return ResponseMsgUnknownError
}
