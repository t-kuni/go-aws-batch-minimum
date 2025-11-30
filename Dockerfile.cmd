# ビルドステージ
FROM golang:1.25.4 AS builder

WORKDIR /app

# Go modulesの初期化とダウンロード
COPY go.mod go.sum* ./
RUN if [ -f go.sum ]; then go mod download; fi

# ソースコードをコピー
COPY . .

# バイナリをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

# 実行ステージ
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# ビルドステージからバイナリをコピー
COPY --from=builder /app/main .

# 実行
CMD ["./main"]