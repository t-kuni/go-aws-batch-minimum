# go-aws-batch-minimum

## 環境変数

| 環境変数名 | 説明 | デフォルト値 | 備考 |
|-----------|------|------------|------|
| APP_WAIT | 処理を待機する時間（秒） | なし | 正の整数を指定。指定がない場合は待機しない |
| APP_RESULT | 処理結果の指定 | なし | "FAIL" を指定するとexitコード1で終了 |
| APP_SF_TASK_TOKEN | Step Functionsタスクトークン | なし | 指定時はStep FunctionsのSendTaskSuccess/SendTaskFailureを呼び出す |

## デプロイ手順

```
AWS_ACCOUNT_ID=XXXXX
IMG_TAG=${AWS_ACCOUNT_ID}.dkr.ecr.ap-northeast-1.amazonaws.com/step-func-example/go-aws-batch-minimum

# ECRにログイン
aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin ${IMG_TAG}

# イメージをビルド
docker build -t ${IMG_TAG} .

# イメージをプッシュ
docker push ${IMG_TAG}

# コンテナを実行
docker run --rm ${IMG_TAG}
```