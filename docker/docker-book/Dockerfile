FROM golang:1.24.1

RUN apt-get update && apt-get install -y curl

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./

RUN go build -o /task_manager

EXPOSE 8080

CMD ["/task_manager"]