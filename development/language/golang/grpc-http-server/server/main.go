package main

import (
	"context"
	"grpc-http-server/app"
	"grpc-http-server/config"
)

func main() {
	application := app.NewApplication(config.Config{
		Host: "0.0.0.0",
		Port: 8080,
	})

	ctx := context.Background()
	if err := application.Start(ctx); err != nil {
		panic(err.Error())
	}
}
