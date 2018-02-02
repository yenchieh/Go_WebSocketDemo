FROM golang:1.9.3
RUN mkdir -p /go/src/github.com/yenchieh/Go_WebSocketDemo
ADD ./build /go/src/github.com/yenchieh/Go_WebSocketDemo/
ADD ./view /go/src/github.com/yenchieh/Go_WebSocketDemo/view
WORKDIR /go/src/github.com/yenchieh/Go_WebSocketDemo
CMD ["/go/src/github.com/yenchieh/Go_WebSocketDemo/main"]

EXPOSE 8443
