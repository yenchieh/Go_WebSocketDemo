# Demo Chat room with WebSocket

This is demo for the [Bowery Golang Meetup](https://www.meetup.com/Bowery-Go/events/246465456/)

## This project is using
- Golang
- [Gin Framework](https://github.com/gin-gonic/gin)

## Setup
- Run it with go. The default server port is 8080
```
go run main.go
```
- Run it with Docker
```
make run-docker
```
- Create Docker image
```
make docker-image
```
- Push docker image. (Update the Makefile to use your own docker repo)
```
make push-image
```