FROM golang:1.18.3-alpine3.15

WORKDIR /app

RUN ls -la

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

CMD ["go", "run", "/app/api/main.go"]