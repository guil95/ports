# Ports

[![build](https://github.com/guil95/ports/actions/workflows/app.yml/badge.svg)](https://github.com/guil95/ports/actions/workflows/app.yml)
[![codecov](https://codecov.io/github/guil95/ports/branch/main/graph/badge.svg?token=712UK1A1YN)](https://codecov.io/github/guil95/ports)

![ports](.github/images/ports.gif)

# Usage

To setup the dependencies you need run the follow command
```shell
docker-compose up
```

After that you'll need run the follow command to run the application
``` 
go run cmd/main.go -file="PATH_RELATIVE/ports.json"
```
Default value for the flag `-file` is `ports.json` and the file need be in the root of your application

To run the tests can you execute that command
```shell
go test -v ./...
```