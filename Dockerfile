# builder
FROM golang:1.17-alpine as builder
MAINTAINER contact@berylleta.dev

RUN mkdir -p /build
WORKDIR /build

COPY go.mod go.sum *.go ./
RUN go mod download
RUN mkdir -p /tmp

RUN GOOS=linux GOARCH=amd64 go build -o iam-checker-linux-amd64 .

# final image
FROM alpine:3.15.0
COPY --from=builder /build/iam-checker-linux-amd64 .

# executable
RUN chmod +x /iam-checker-linux-amd64

EXPOSE 80

ENTRYPOINT [ "./iam-checker-linux-amd64" ]