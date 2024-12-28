# example golang

## stack

- docker
- docker-compose
- go gin
- mongodb
- swagger
- jwt
- air
- rate limit
- websocket
- upload file

## setup

- install go and setup path
- install docker and docker-compose
- install and setup air $go install github.com/air-verse/air@latest
- install swag $go install github.com/swaggo/swag/cmd/swag@latest

### create file .air.toml

```
root = "."
tmp_dir = "tmp"

[build]
cmd = "go build -o ./tmp/main.exe ./cmd/app"
bin = "tmp/main.exe"
full_bin = "./tmp/main.exe"
include_ext = ["go"]
exclude_dir = ["tmp", "mongodb_data"]
delay = 1000

[screen]
clear_on_rebuild = true
```

- init project $go mod init example-go-project
- init package $go mod tidy
- cp .env.example .env
- init swagger $swag init -g cmd/app/main.go
- init cors $go get github.com/gin-contrib/cors
- build $go build cmd/app/main.go

## how to use

- run $docker-compose up -d --build (init project or db)
- run app $go run cmd/app/main.go or use $air (air is build and compiler follow code change)

## run test

- $go test ./internal/handlers/user/test
- $go test -race ./internal/handlers/user/test -v -cover
- $go test -race ./internal/handlers/user/test -v -coverprofile=coverage.out && go tool cover -html=coverage.out

## Feature

- use swagger [x]
- use ratelimit [x]
- use jwt [x]
- use mongodb [x]
- use auth [x]
- use call external api [x]
- use upload and read file [x]
- use docker [x]
- set pattern code [x]
- unit test [x]
- restful api [x]
- relation db [x]
- permission roles [x]
- pagination [x]
