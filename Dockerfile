FROM golang:1.17.5-alpine3.15

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o json-fixer

EXPOSE 9900

CMD sh -c "/app/json-fixer"