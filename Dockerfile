FROM golang:alpine AS builder

RUN apk --update add git
RUN mkdir -p /go/src/hello_world
WORKDIR /go/src/hello_world
COPY *.go .
RUN go get -v && go build

FROM alpine
COPY --from=builder /go/src/hello_world/hello_world /hello_world

EXPOSE 9000

CMD ["/hello_world"]
