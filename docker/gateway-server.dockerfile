# build stage
FROM golang:1.26.2-alpine3.22 AS builder

WORKDIR /build
COPY /backend/gateway-server/go.mod ./
COPY /backend/gateway-server/cmd  ./cmd 
COPY  /backend/gateway-server/internal ./internal


RUN go mod tidy

RUN go build  -o ./gateway-server ./cmd

# final stage
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /build/gateway-server .
COPY /backend/gateway-server/scripts/start.sh .

COPY /configs/gateway-server/dev.yaml ./
ENV CONFIG_PATH=./dev.yaml

RUN chmod +x ./start.sh
CMD ["sh", "./start.sh"]