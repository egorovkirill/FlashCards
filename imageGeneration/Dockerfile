FROM golang:latest

RUN go version
ENV GOPATH=/

WORKDIR /image
COPY .. /image


EXPOSE 8000

RUN go mod download
RUN go build -o main ./cmd/main.go
CMD ["/image/main"]
