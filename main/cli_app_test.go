package main

import (
	"context"
	"testing"
)

func TestConnectTogether(t *testing.T) {
	// Given
	ctx := context.Background()
	err := RunCLIApp(ctx)
	if err != nil {
		t.Errorf("error shouldn't happen here, \nexpected: %v, \nactual: %v", nil, err)
		return
	}
}
