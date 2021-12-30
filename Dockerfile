# syntax=docker/dockerfile:1
##
## Build
##
FROM golang:1.17 As builder

ENV GOPROXY=https://goproxy.cn

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY configs/config.yaml ./configs/

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o syncer ./cmd/syncer

##
## Deploy
##
FROM alpine:3.14

COPY --from=builder /app/syncer /syncer
COPY --from=builder /app/configs/config.yaml /configs/config.yaml
RUN chmod +x /syncer

CMD ["/syncer"]