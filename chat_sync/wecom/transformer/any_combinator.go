package transformer

type AnyCombinator struct {
	fetchers []NameFetcher
}

func NewAnyCombinator(fetchers []NameFetcher) *AnyCombinator {
	if len(fetchers) == 0 {
		panic("fetchers can't be empty")
	}

	return &AnyCombinator{
		fetchers: fetchers,
	}
}

func (c *AnyCombinator) FetchName(id string) (name string, err error) {
	for _, fetcher := range c.fetchers {
		name, err = fetcher.FetchName(id)
		if err == nil {
			return name, nil
		}
	}

	return "", err
}
