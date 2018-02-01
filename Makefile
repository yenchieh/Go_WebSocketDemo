.PHONY: deps dev build clean push-image

deps:
	dep ensure

build: clean deps
	GOOS=linux go build -o ./build/main ./*.go

clean:
	rm -rf build/*

build-image:
	docker build --rm -t yenchieh/websocket_demo:latest .

push-image: build build-image
	docker push yenchieh/websocket_demo:latest

run:
	go run main.go