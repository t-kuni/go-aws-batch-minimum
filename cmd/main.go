package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/t-kuni/go-aws-batch-minimum/app"
)

func main() {
	app := app.NewApp()

	waitSecondsStr := os.Getenv("APP_WAIT")
	if waitSecondsStr != "" {
		waitSeconds, err := strconv.Atoi(waitSecondsStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: APP_WAIT must be a valid integer: %v\n", err)
			os.Exit(1)
		}
		app.SetWaitSeconds(waitSeconds)
	}

	result := os.Getenv("APP_RESULT")
	if result != "" {
		app.SetResult(result)
	}
	
	taskToken := os.Getenv("APP_SF_TASK_TOKEN")
	if taskToken != "" {
		app.SetTaskToken(taskToken)
	}

	itemsCountStr := os.Getenv("APP_RESULT_ITEMS_COUNT")
	if itemsCountStr != "" {
		itemsCount, err := strconv.Atoi(itemsCountStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: APP_RESULT_ITEMS_COUNT must be a valid integer: %v\n", err)
			os.Exit(1)
		}
		app.SetResultItemsCount(itemsCount)
	}

	app.Exec()
}
