package name_fetcher

import (
	"github.com/stretchr/testify/assert"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom/transformer"
	"testing"
)

func TestAnyCombinator_FetchName_ZeroFetcher(t *testing.T) {
	// if given zero fetcher, the constructor should panic.
	// use assert library to check if the constructor panics.

	assert.Panics(t, func() {
		NewAnyCombinator(nil)
	})

	assert.Panics(t, func() {
		NewAnyCombinator([]transformer.NameFetcher{})
	})
}
