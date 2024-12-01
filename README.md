### example golang

## stack

- docker
- docker-compose
- go gin
- mongodb
- swagger
- jwt

## setup

- init project $go mod init example-go-project
- init package $go mod tidy
- cp .env.example .env
- init swagger $swag init -g cmd/app/main.go
- update api swagger $swag init -g cmd/app/main.go -d ./

## how to use

- run $docker-compose up -d --build (init project or db)
- run app $go run cmd/app/main.go
