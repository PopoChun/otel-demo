FROM golang:1.19.2 as builder
WORKDIR /go/src/otel-bar
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /go/src/otel-bar/main /app/main
ENTRYPOINT ["/app/main"]