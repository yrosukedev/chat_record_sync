package paginated_reader

type PageToken struct {
	Value int64
}

func NewPageToken(value int64) *PageToken {
	return &PageToken{
		Value: value,
	}
}
