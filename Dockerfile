FROM golang:1.22-alpine

WORKDIR /app

COPY ./database ./database
COPY ./handler ./handler
COPY ./main ./main
COPY ./public ./public
COPY ./helper ./helper
COPY go.mod .
RUN go mod tidy

RUN go build -o -tags DE /poliSim /app/main

EXPOSE 8080

CMD ["/poliSim"]