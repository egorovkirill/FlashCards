FROM golang:latest

RUN go version
ENV GOPATH=/

WORKDIR /image
COPY .. /image

RUN go mod download
RUN go build -o main ./cmd/main.go
CMD ["/image/main"]
