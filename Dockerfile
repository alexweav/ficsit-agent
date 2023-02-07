FROM golang:1.20-alpine as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . .

RUN go build -v -o ficsit-agent cmd/main.go

FROM alpine:3.17.1

COPY --from=builder /app/ficsit-agent /app/ficsit-agent

ENTRYPOINT ["/app/ficsit-agent"]
