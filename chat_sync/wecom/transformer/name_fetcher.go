package transformer

type NameFetcher interface {
	FetchName(id string) (name string, err error)
}
