package paginated_reader

type PageToken struct {
	Value uint64
}

func NewPageToken(value uint64) *PageToken {
	return &PageToken{
		Value: value,
	}
}
