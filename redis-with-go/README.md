# Redis with golang

## Install redis from docker
Pull the image from docker
```
docker pull redis
```

Start the redis container
```
docker run --name <redis-name> -p 6379:6379 -d redis
```

List the running containers
```
docker ps
```

Stop the redis container
```
docker stop my-redis-container
```

## Create golang project
Init go project
```
go mod init <project-name>
```

Download redis package
```
go get github.com/go-redis/redis/v8
```

Build the go project
```
go build
```

Run
```
go run main.go
```