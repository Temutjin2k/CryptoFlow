FROM golang:1.24

WORKDIR /app

COPY . .

RUN go build -o main cmd/main.go

EXPOSE 8080

CMD ["./cmd/main.go"]