FROM golang:1.25.1

WORKDIR /app

COPY . .

#For local run use RUN go build -o main ./lastthirdapp/main.go
RUN go build -o main ./lastthirdapp/main.go

EXPOSE 8080

CMD ["./main"]