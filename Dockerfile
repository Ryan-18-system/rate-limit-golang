FROM golang:latest as builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o server ./cmd/server

# Final image
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/server .

CMD ["./server"]