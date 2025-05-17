ARG GO_VERSION=1.24.3
ARG ALPINE_VERSION=3.16

### Build binary
FROM golang:${GO_VERSION} as build-binary
COPY . /project
WORKDIR /project
RUN go mod vendor \
    && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 GO111MODULE=on go build \
        -v \
        -mod vendor \
        -o /project/bin/todo_app \
            /project/cmd/todo_app

### Image
FROM alpine:${ALPINE_VERSION} as image
COPY --from=build-binary /project/bin/todo_app /usr/local/bin/todo_app
RUN chmod +x /usr/local/bin/todo_app
ENTRYPOINT ["todo_app"]
