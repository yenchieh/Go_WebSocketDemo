.PHONY: deps dev build clean push-image run run-docker

deps:
	go mod tidy

build: clean deps
	GOOS=linux go build -o ./build/main ./*.go

clean:
	rm -rf build/*

build-image:
	docker build --rm -t websocket_demo:latest .

push-image: build build-image
	docker push websocket_demo:latest

run:
	go run main.go

run-docker: build build-image
	docker run --rm -d -p 8080:8080 websocket_demo:latest
