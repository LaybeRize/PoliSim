FROM golang:1.22-alpine

WORKDIR /app

COPY ./handler ./handler
COPY ./database ./database
COPY ./templates ./templates
COPY ./main ./main
COPY ./public ./public
COPY ./resources ./resources
COPY go.mod .
RUN go mod tidy

RUN go build -o /poliSim /app/main

EXPOSE 8080

CMD ["/poliSim"]