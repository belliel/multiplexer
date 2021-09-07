FROM golang:1.16.7-alpine3.14

WORKDIR /app

COPY ./ ./

ENV LISTEN_ADDR=":8080"
ENV DEBUG=TRUE
ENV CONNECTIONS_LIMIT=100

EXPOSE 8080:8080

RUN go build -o main cmd/app/main.go

CMD ./main
