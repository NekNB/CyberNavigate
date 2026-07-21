# build stage
FROM golang:1.26.2-alpine3.22 AS builder

WORKDIR /build
COPY /backend/simulator-service/go.mod ./
COPY /backend/simulator-service/cmd  ./cmd 
COPY  /backend/simulator-service/internal ./internal


RUN go mod tidy

RUN go build  -o ./simulator-service ./cmd

# # final stage
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /build/simulator-service .
COPY /backend/simulator-service/scripts/start.sh .

COPY /configs/simulator-service/dev.yaml ./
COPY /secrets/certs/simulator-service ./certs

ENV CONFIG_PATH=./dev.yaml

RUN chmod +x ./start.sh
RUN sed -i 's/\r$//'  ./start.sh
CMD ["sh", "./start.sh"]