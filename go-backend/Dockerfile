FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd
RUN go build -o /app/main .

WORKDIR /app
CMD ["./main"]
