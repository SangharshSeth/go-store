FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -v -o app main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/app .
COPY --from=builder /app/templates ./templates
EXPOSE 5000

CMD ["./app"]
