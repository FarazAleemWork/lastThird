FROM golang:1.25.1-alpine

WORKDIR /app/lastThirdApp

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o bin/lastthird ./lastThirdApp

EXPOSE 8080

ENTRYPOINT [ "/app/lastThirdApp/bin/lastthird" ]