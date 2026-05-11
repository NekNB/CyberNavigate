# build stage
FROM golang:1.26.2-trixie AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod tidy


COPY ./cmd ./internal ./
RUN go build ./cmd -o user-service

# final stage
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /build/user-service .
COPY ./scripts/start.sh .
ENTRYPOINT [ "./scripts/" ]