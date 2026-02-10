FROM golang:1.25 AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o url-shortener-go ./cmd/main.go

FROM alpine:3.20
WORKDIR /app
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
COPY --from=builder /app/url-shortener-go /app/url-shortener-go
USER appuser

EXPOSE 8080

ENV HOST=0.0.0.0
ENV PORT=8080

ENTRYPOINT ["/app/url-shortener-go"]