package main

import (
	"context"
	"github.com/yrosukedev/chat_record_sync/http_app"
)

func main() {
	ctx := context.Background()
	app := http_app.NewHTTPApp(ctx)
	app.Run(ctx)
}
