FROM golang:1.11-alpine3.8 AS builder

RUN apk --update add git
WORKDIR /src/hello_world
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build

FROM alpine:3.8
COPY --from=builder /src/hello_world/hello-world /hello-world

EXPOSE 9000

CMD ["/hello-world"]
