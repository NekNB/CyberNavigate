# build stage
FROM golang:1.26.2-alpine3.22 AS builder

WORKDIR /build
COPY /backend/article-service/go.mod /backend/article-service/go.sum ./

RUN go mod download -x

COPY /backend/article-service/cmd  ./cmd 
COPY  /backend/article-service/internal ./internal

RUN go build  -o ./article-service ./cmd

# final stage
FROM alpine:latest

WORKDIR /root/
COPY --from=builder /build/article-service .
COPY /backend/article-service/scripts/start.sh .

COPY /configs/article-service/dev.yaml ./
COPY /secrets/certs/article-service ./certs

ENV CONFIG_PATH=./dev.yaml

RUN chmod +x ./start.sh
RUN sed -i 's/\r$//'  ./start.sh
CMD ["sh", "./start.sh"]