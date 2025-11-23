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