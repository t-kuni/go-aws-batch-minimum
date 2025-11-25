package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
)

// TaskOutput は Step Functionsに返すペイロードの構造体です
type TaskOutput struct {
	Status string `json:"status"`
	Result string `json:"result"`
	Items  []int  `json:"items"`
}

func main() {
	fmt.Println("Hello, World!")

	// 全ての環境変数を標準出力に出力
	fmt.Println("\n=== Environment Variables ===")
	for _, env := range os.Environ() {
		fmt.Println(env)
	}
	fmt.Println("=============================\n")

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
	taskToken := os.Getenv("APP_SF_TASK_TOKEN")
	itemsCountStr := os.Getenv("APP_RESULT_ITEMS_COUNT")

	// Step Functionsのタスクトークンが指定されている場合
	if taskToken != "" {
		ctx := context.Background()

		// AWS設定とStep Functionsクライアントの初期化
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			log.Fatalf("unable to load AWS config, %v", err)
		}

		sfnClient := sfn.NewFromConfig(cfg)

		if result == "FAIL" {
			// 失敗の場合: SendTaskFailureを呼び出す
			errorMessage := "APP_RESULT is FAIL. Task execution failed."
			log.Printf("Task failed. Sending SendTaskFailure for token: %s", taskToken)

			_, err := sfnClient.SendTaskFailure(ctx, &sfn.SendTaskFailureInput{
				TaskToken: &taskToken,
				Error:     stringPtr("TaskExecutionFailed"),
				Cause:     stringPtr(errorMessage),
			})

			if err != nil {
				log.Fatalf("Failed to send task failure: %v", err)
			}
			log.Println("SendTaskFailure completed successfully.")
			os.Exit(1)
		} else {
			// 成功の場合: SendTaskSuccessを呼び出す
			// APP_RESULT_ITEMS_COUNTに基づいてitemsを生成
			items := []int{}
			if itemsCountStr != "" {
				itemsCount, err := strconv.Atoi(itemsCountStr)
				if err != nil {
					log.Fatalf("Error: APP_RESULT_ITEMS_COUNT must be a valid integer: %v", err)
				}
				if itemsCount > 0 {
					items = make([]int, itemsCount)
					for i := 0; i < itemsCount; i++ {
						items[i] = i + 1
					}
				}
			}

			output := TaskOutput{
				Status: "SUCCESS",
				Result: fmt.Sprintf("Task completed at %s", time.Now().Format(time.RFC3339)),
				Items:  items,
			}

			outputJSON, err := json.Marshal(output)
			if err != nil {
				log.Fatalf("failed to marshal output JSON: %v", err)
			}

			log.Printf("Task succeeded. Sending SendTaskSuccess for token: %s", taskToken)
			_, err = sfnClient.SendTaskSuccess(ctx, &sfn.SendTaskSuccessInput{
				TaskToken: &taskToken,
				Output:    stringPtr(string(outputJSON)),
			})

			if err != nil {
				log.Fatalf("Failed to send task success: %v", err)
			}
			log.Println("SendTaskSuccess completed successfully.")
		}
	} else {
		// タスクトークンが指定されていない場合は従来の動作
		if result == "FAIL" {
			fmt.Println("APP_RESULT is FAIL. Exiting with code 1.")
			os.Exit(1)
		}
	}

	fmt.Println("Process completed successfully.")
}

// stringPtr は文字列のポインタを返すヘルパー関数
func stringPtr(s string) *string {
	return &s
}
