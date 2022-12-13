# Ports

![ports](.github/images/ports.gif)

# Usage

To setup the database you need run the follow command
```shell
docker-compose up
```

After that you need run the follow command to run the application
``` 
go run cmd/main.go -file="PATH_RELATIVE/ports.json"
```

The default value to -file is `ports.json` and the file need be in the root of your application