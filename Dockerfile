FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o main ./cmd/app

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app
COPY --from=builder /app/main .

CMD ["./main"]
