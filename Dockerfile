# syntax=docker/dockerfile:1
##
## Build
##
FROM golang:1.17 As builder

ENV GOPROXY=https://goproxy.io,direct

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o syncer ./cmd/syncer

##
## Deploy
##
FROM alpine:3.14

COPY --from=builder /app/syncer /syncer
RUN chmod +x /syncer

CMD ["/syncer"]