# MongoDB with golang

## Install mongodb from docker
Pull the image from docker
```
docker pull mongo
```

Start the container
```
docker run --name <name> -p 27017:27017 -d mongo
```

List the running containers
```
docker ps
```

Stop the container
```
docker stop <mongo-name>
```

## Create golang project
Init go project
```
go mod init <project-name>
```

Download mongodb package
```
go get go.mongodb.org/mongo-driver/mongo
```
```
go get go.mongodb.org/mongo-driver/bson
```

Build the go project
```
go build
```

Run
```
go run main.go
```