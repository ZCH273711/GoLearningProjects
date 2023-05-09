FROM golang:latest
LABEL authors="admin"

ENV GO111MODULE=on \
    CGO_ENABLE=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build
COPY . .
RUN go mod tidy && cd ./data && go mod tidy && cd ../ && go build -o app
EXPOSE 8080
ENTRYPOINT ["./app"]
