#first stage - builder
FROM golang:1.17.13-alpine3.16 as builder
WORKDIR /build
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY  . .
ENV GO111MODULE=on
RUN CGO_ENABLED=0 GOOS=linux go build -o majapahit-service

#second stage
FROM alpine:3.17.0
WORKDIR /root/
COPY --from=builder /build .
CMD ["./majapahit-service"]
