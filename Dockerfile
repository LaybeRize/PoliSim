FROM golang:1.21-alpine

WORKDIR /app

COPY ./data ./data
COPY ./helper ./helper
COPY ./html ./html
COPY ./main ./main
COPY ./public ./public
COPY ./resources ./resources
COPY go.mod .
RUN go mod tidy

RUN go build -o /poliSim /app/main

EXPOSE 8080

CMD ["/poliSim"]