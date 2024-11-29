# Build Stage
FROM golang:1.17-alpine AS builder

WORKDIR /app

# 필요한 패키지 설치
RUN apk add --no-cache git

# 의존성 설치 및 빌드
COPY go.mod go.sum ./
RUN go mod download
COPY consumer.go ./
RUN go build -o consumer consumer.go

# Runtime Stage
FROM alpine:3.16

WORKDIR /app

# 빌드된 실행 파일 복사
COPY --from=builder /app/consumer /app/

CMD ["./consumer"]
