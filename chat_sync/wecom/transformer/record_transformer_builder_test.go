package transformer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRecordTransformerBuilder_Build_OpenAPIServiceMayBeNil(t *testing.T) {
	// if the openAPIService is nil, the build should not panic
	// we should use assert library to check the panic

	// Given
	builder := NewRecordTransformerBuilder(nil)

	// Then
	assert.NotPanics(t, func() {
		// When
		builder.Build()
	})
}

