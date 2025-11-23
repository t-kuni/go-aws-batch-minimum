package main

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	fmt.Println("Hello, World!")

	// APP_WAIT環境変数を取得（単位：秒）
	waitSecondsStr := os.Getenv("APP_WAIT")
	if waitSecondsStr != "" {
		waitSeconds, err := strconv.Atoi(waitSecondsStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: APP_WAIT must be a valid integer: %v\n", err)
			os.Exit(1)
		}
		if waitSeconds > 0 {
			fmt.Printf("Waiting for %d seconds...\n", waitSeconds)
			time.Sleep(time.Duration(waitSeconds) * time.Second)
			fmt.Println("Wait completed.")
		}
	}

	// APP_RESULT環境変数を確認
	result := os.Getenv("APP_RESULT")
	if result == "FAIL" {
		fmt.Println("APP_RESULT is FAIL. Exiting with code 1.")
		os.Exit(1)
	}

	fmt.Println("Process completed successfully.")
}
