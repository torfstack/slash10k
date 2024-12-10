FROM golang:1.23.4-alpine3.21 AS builder

RUN mkdir /opt/slash10k
WORKDIR /opt/slash10k

COPY go.mod go.sum ./
RUN go mod download

COPY cmd cmd
COPY pkg pkg
COPY sql sql

RUN CGO_ENABLED=0 GOOS=linux go build -o bin/slash10k cmd/bot/main.go

FROM alpine:3.21.0

RUN mkdir /opt/slash10k
WORKDIR /opt/slash10k

COPY --from=builder /opt/slash10k/bin/slash10k slash10k

CMD ["./slash10k"]
