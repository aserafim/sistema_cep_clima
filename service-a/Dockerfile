FROM golang:1.23-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o main .

FROM alpine:3.18
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]