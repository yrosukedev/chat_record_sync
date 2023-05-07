package name_fetcher

import (
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/yrosukedev/chat_record_sync/chat_sync/wecom/transformer"
	"io"
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

func TestAnyCombinator_FetchName_OneFetcher(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	fetcher := NewMockNameFetcher(ctrl)
	fetcher.EXPECT().FetchName("123").Return("haary", nil).Times(1)
	combinator := NewAnyCombinator([]transformer.NameFetcher{fetcher})

	// When
	name, err := combinator.FetchName("123")

	// Then
	if assert.NoError(t, err) {
		assert.Equal(t, "haary", name)
	}
}

func TestAnyCombinator_FetchName_OneFetcher_Error(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	fetcher := NewMockNameFetcher(ctrl)
	fetcher.EXPECT().FetchName("123").Return("", io.EOF).Times(1)
	combinator := NewAnyCombinator([]transformer.NameFetcher{fetcher})

	// When
	_, err := combinator.FetchName("123")

	// Then
	assert.Error(t, err)
}

func TestAnyCombinator_FetchName_MultipleFetchers_FirstFetcherFails(t *testing.T) {
	// Given
	ctrl := gomock.NewController(t)
	fetcher1 := NewMockNameFetcher(ctrl)
	fetcher1.EXPECT().FetchName("123").Return("", io.EOF).Times(1)
	fetcher2 := NewMockNameFetcher(ctrl)
	fetcher2.EXPECT().FetchName("123").Return("haary", nil).Times(1)
	fetcher3 := NewMockNameFetcher(ctrl)
	fetcher3.EXPECT().FetchName("123").Return("marry", nil).Times(0)
	combinator := NewAnyCombinator([]transformer.NameFetcher{fetcher1, fetcher2, fetcher3})

	// When
	name, err := combinator.FetchName("123")

	// Then
	if assert.NoError(t, err) {
		assert.Equal(t, "haary", name)
	}
}
