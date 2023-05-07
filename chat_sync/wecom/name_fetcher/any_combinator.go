package name_fetcher

import (
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom/transformer"
)

type AnyCombinator struct {
	fetchers []transformer.NameFetcher
}

func NewAnyCombinator(fetchers []transformer.NameFetcher) *AnyCombinator {
	if len(fetchers) == 0 {
		panic("fetchers can't be empty")
	}

	return &AnyCombinator{
		fetchers: fetchers,
	}
}

func (c *AnyCombinator) FetchName(id string) (name string, err error) {
	// TODO: implement me
	panic("implement me")
}
