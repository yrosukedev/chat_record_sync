package paginated_reader

type ChatRecordPaginationStorage interface {
	Get() (pageToken PageToken, err error)
	Set(pageToken PageToken) error
}
