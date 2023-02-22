FROM golang:latest

RUN go version
ENV GOPATH=/

WORKDIR /app
COPY .. /app

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

RUN chmod +x cold-start-postgres.sh

RUN go mod download
RUN go build -o main ./cmd/main.go
CMD ["/app/main"]
