FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/go-horses-racing-by-cobra ./cmd

FROM scratch

COPY --from=builder /app/go-horses-racing-by-cobra /go-horses-racing-by-cobra

ENTRYPOINT ["/go-horses-racing-by-cobra"]