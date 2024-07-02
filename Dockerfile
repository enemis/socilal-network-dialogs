FROM golang:1.21

WORKDIR /usr/src/app
COPY . .
ENV GOPROXY=https://goproxy.io,direct
RUN go mod download
RUN go mod verify
RUN go build -v -o /usr/local/bin/app main.go
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
CMD ["app"]