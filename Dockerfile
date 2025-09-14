FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV TODO_PORT=7540

RUN GOOS=linux GOARCH=amd64 go build -o /finalTask

EXPOSE $TODO_PORT

CMD ["/finalTask"] 