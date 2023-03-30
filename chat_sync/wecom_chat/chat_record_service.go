package wecom_chat

type ChatRecordService interface {
	Read(seq uint64, pageSize uint64) (records []*WeComChatRecord, err error)
}
