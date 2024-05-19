FROM golang:1.19-alpine

WORKDIR /app

RUN apk add --no-cache docker

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

EXPOSE 3000

CMD ["./main"]
