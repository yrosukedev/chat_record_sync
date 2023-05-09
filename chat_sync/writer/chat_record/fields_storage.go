package chat_record

type FieldsStorage interface {
	Write(fields map[string]interface{}, requestUUID string) error
}
