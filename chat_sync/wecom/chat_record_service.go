package wecom

type ChatRecordService interface {
	Read(seq uint64, pageSize uint64) (records []*ChatRecord, err error)
}
