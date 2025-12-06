package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/t-kuni/go-aws-batch-minimum/app"
)

func handler(ctx context.Context) (string, error) {
	app.Exec()
	return "Hello, Lambda!", nil
}

func main() {
	lambda.Start(handler)
}
