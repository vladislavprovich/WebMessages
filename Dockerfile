FROM golang:1.23.6 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o messenger ./cmd


FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/messenger .

EXPOSE 8080

CMD ["./messenger"]
