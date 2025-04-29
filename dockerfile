FROM golang:1.24.2 AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY .env .env

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o apiserver .

FROM scratch

COPY --from=builder /build/apiserver /
COPY --from=builder /build/.env /

ENTRYPOINT ["/apiserver"]