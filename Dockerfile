FROM golang:1.24rc2-alpine

WORKDIR /app

COPY ./database ./database
COPY ./handler ./handler
COPY ./main ./main
COPY ./public ./public
COPY ./helper ./helper
COPY ./localisation ./localisation
COPY go.mod .
RUN go mod tidy

RUN go build -tags DE -o /poliSim /app/main

EXPOSE 8080

CMD ["/poliSim"]