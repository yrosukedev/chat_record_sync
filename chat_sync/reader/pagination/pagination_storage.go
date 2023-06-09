package pagination

type PaginationStorage interface {
	// Get return the page token.
	// Return nil if the page token is never stored before.
	Get() (pageToken *PageToken, err error)

	Set(pageToken *PageToken) error
}
