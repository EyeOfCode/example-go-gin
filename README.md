# example golang

## stack

- docker
- docker-compose
- go gin
- mongodb
- swagger
- jwt
- dogo
- ratelimit

## setup

- install and setup dogo $go get github.com/liudng/dogo

### crate file dogo.json

```
{
    "WorkingDir": "{GOPATH}/src/github.com/liudng/dogo/example",
    "SourceDir": [
        "{GOPATH}/src/github.com/liudng/dogo/example"
    ],
    "SourceExt": [".c", ".cpp", ".go", ".h"],
    "BuildCmd": "go build github.com/liudng/dogo/example",
    "RunCmd": "example.exe",
    "Decreasing": 1
}
```

- init project $go mod init example-go-project
- init package $go mod tidy
- cp .env.example .env
- init swagger $swag init -g cmd/app/main.go
- update api swagger $swag init -g cmd/app/main.go -d ./
- build $go build cmd/app/main.go

## how to use

- run $docker-compose up -d --build (init project or db)
- run app $go run cmd/app/main.go or use $dogo (dogo is build and compiler follow code change)

## TODO

- use swagger [x]
- use ratelimit [x]
- use jwt [x]
- use mongodb [x]
- use auth [x]
- use call external api [ ]
- use upload file [ ]
- use docker [x]
- set pattern code [ ]
