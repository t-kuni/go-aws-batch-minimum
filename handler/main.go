package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/t-kuni/go-aws-batch-minimum/app"
)

type StepFunctionInput struct {
	WaitSeconds      int     `json:"wait_seconds"`
	Result           string  `json:"result"`
	TaskToken        string  `json:"task_token"`
	ResultItemsCount int     `json:"result_items_count"`
	FailPercent      float64 `json:"fail_percent"`
}

func handler(ctx context.Context, event StepFunctionInput) (string, error) {
	app := app.NewApp()

	if event.WaitSeconds > 0 {
		app.SetWaitSeconds(event.WaitSeconds)
	}

	if event.Result != "" {
		app.SetResult(event.Result)
	}

	if event.TaskToken != "" {
		app.SetTaskToken(event.TaskToken)
	}

	if event.ResultItemsCount > 0 {
		app.SetResultItemsCount(event.ResultItemsCount)
	}

	if event.FailPercent > 0 {
		app.SetFailPercent(event.FailPercent)
	}

	app.Exec()
	return "Hello, Lambda!", nil
}

func main() {
	lambda.Start(handler)
}
