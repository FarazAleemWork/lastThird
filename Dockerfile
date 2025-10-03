FROM golang:1.25.1

WORKDIR /app

COPY . .

RUN go build -o main ./lastthirdapp/main.go

EXPOSE 8080

CMD ["./main"]