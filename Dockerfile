FROM golang:1.21-alpine

WORKDIR /app

COPY ./componentHelper ./componentHelper
COPY ./database ./database
COPY ./dataExtraction ./dataExtraction
COPY ./dataValidation ./dataValidation
COPY ./htmlComposition ./htmlComposition
COPY ./htmlServer ./htmlServer
COPY ./main ./main
COPY ./public ./public
COPY ./resources ./resources
COPY go.mod .
RUN go mod tidy

RUN go build -o /mbundestag /app/main

EXPOSE 8080

CMD ["/mbundestag"]