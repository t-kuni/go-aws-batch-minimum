package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/t-kuni/go-aws-batch-minimum/app"
)

type StepFunctionInput struct {
	WaitSeconds      int    `json:"wait_seconds"`
	Result           string `json:"result"`
	TaskToken        string `json:"task_token"`
	ResultItemsCount int    `json:"result_items_count"`
}

func handler(ctx context.Context, event StepFunctionInput) (string, error) {
	app := app.NewApp()
	app.Exec()
	return "Hello, Lambda!", nil
}

func main() {
	lambda.Start(handler)
}
