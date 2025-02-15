FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/main ./cmd/main.go

FROM alpine:latest

COPY --from=builder /bin/main /bin/main

COPY css/ css/
COPY images/ images/
COPY views/ views/

EXPOSE 8080

CMD ["/bin/main"]
