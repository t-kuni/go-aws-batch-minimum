package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
)

type app struct {
	waitSeconds      int
	result           string
	taskToken        string
	resultItemsCount int
	failPercent      float64
}

// TaskOutput は Step Functionsに返すペイロードの構造体です
type TaskOutput struct {
	Status string `json:"status"`
	Result string `json:"result"`
	Items  []int  `json:"items"`
}

func NewApp() *app {
	return &app{}
}

func (a *app) SetWaitSeconds(waitSeconds int) *app {
	a.waitSeconds = waitSeconds
	return a
}

func (a *app) SetResult(result string) *app {
	a.result = result
	return a
}

func (a *app) SetTaskToken(taskToken string) *app {
	a.taskToken = taskToken
	return a
}

func (a *app) SetResultItemsCount(resultItemsCount int) *app {
	a.resultItemsCount = resultItemsCount
	return a
}

func (a *app) SetFailPercent(failPercent float64) *app {
	a.failPercent = failPercent
	return a
}

func (a *app) Exec() {
	fmt.Println("Hello, World!")

	// アプリケーションパラメータを標準出力に出力
	fmt.Println("\n=== Application Parameters ===")
	fmt.Printf("waitSeconds:      %d\n", a.waitSeconds)
	fmt.Printf("result:           %s\n", a.result)
	fmt.Printf("taskToken:        %s\n", a.taskToken)
	fmt.Printf("resultItemsCount: %d\n", a.resultItemsCount)
	fmt.Printf("failPercent:      %.2f\n", a.failPercent)
	fmt.Println("===============================\n")

	// 全ての環境変数を標準出力に出力
	fmt.Println("=== Environment Variables ===")
	for _, env := range os.Environ() {
		fmt.Println(env)
	}
	fmt.Println("=============================\n")

	if a.waitSeconds > 0 {
		fmt.Printf("Waiting for %d seconds...\n", a.waitSeconds)
		time.Sleep(time.Duration(a.waitSeconds) * time.Second)
		fmt.Println("Wait completed.")
	}

	// failPercentに基づいて確率的に失敗させる
	if a.failPercent > 0 {
		randomValue := rand.Float64()
		fmt.Printf("failPercent: %.2f, randomValue: %.2f\n", a.failPercent, randomValue)
		if randomValue < a.failPercent {
			fmt.Println("Random failure triggered based on failPercent")
			a.result = "FAIL"
		}
	}

	// Step Functionsのタスクトークンが指定されている場合
	if a.taskToken != "" {
		ctx := context.Background()

		// AWS設定とStep Functionsクライアントの初期化
		cfg, err := config.LoadDefaultConfig(ctx)
		if err != nil {
			log.Fatalf("unable to load AWS config, %v", err)
		}

		sfnClient := sfn.NewFromConfig(cfg)

		if a.result == "FAIL" {
			// 失敗の場合: SendTaskFailureを呼び出す
			errorMessage := "APP_RESULT is FAIL. Task execution failed."
			log.Printf("Task failed. Sending SendTaskFailure for token: %s", a.taskToken)

			_, err := sfnClient.SendTaskFailure(ctx, &sfn.SendTaskFailureInput{
				TaskToken: &a.taskToken,
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
			items := make([]int, a.resultItemsCount)
			for i := 0; i < a.resultItemsCount; i++ {
				items[i] = i + 1
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

			log.Printf("Task succeeded. Sending SendTaskSuccess for token: %s", a.taskToken)
			_, err = sfnClient.SendTaskSuccess(ctx, &sfn.SendTaskSuccessInput{
				TaskToken: &a.taskToken,
				Output:    stringPtr(string(outputJSON)),
			})

			if err != nil {
				log.Fatalf("Failed to send task success: %v", err)
			}
			log.Println("SendTaskSuccess completed successfully.")
		}
	} else {
		// タスクトークンが指定されていない場合は従来の動作
		if a.result == "FAIL" {
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
