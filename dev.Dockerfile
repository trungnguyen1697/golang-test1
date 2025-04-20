# DO NOT use in production!
FROM golang:1.24-alpine

# install file watcher
RUN apk add make && go get github.com/githubnemo/CompileDaemon && go get github.com/swaggo/swag/cmd/swag

# change working dir
WORKDIR /app

# Copy go module files and download dependencies
COPY go.* ./
RUN go mod download

ENTRYPOINT ["CompileDaemon", "-exclude-dir=.git", "-exclude-dir=docs", "-build=make build", "-command=./build/golang-test1", "-graceful-kill=true"]
