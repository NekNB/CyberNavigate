# build stage
FROM golang:1.26.2-alpine3.22 AS builder

WORKDIR /build
COPY /backend/user-service/go.mod  /backend/user-service/go.sum ./
RUN go mod download -x

COPY /backend/user-service/cmd  ./cmd 
COPY  /backend/user-service/internal ./internal



RUN go build  -o ./user-service ./cmd

# final stage
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /build/user-service .
COPY /backend/user-service/scripts/start.sh .

COPY /configs/user-service/dev.yaml ./
COPY /secrets/keys ./keys
COPY /secrets/certs/user-service ./certs

ENV CONFIG_PATH=./dev.yaml

RUN chmod +x ./start.sh
RUN sed -i 's/\r$//'  ./start.sh
CMD ["sh", "./start.sh"]