# Demo Chat room with WebSocket

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
- Run the docker image
```
docker pull ${repoName}/websocket_demo:latest
docker run -d -e ADDRESS='${Your Server IP}' -e PORT=':8443' -p 8443:8443 ${repoName}/websocket_demo:latest
```
