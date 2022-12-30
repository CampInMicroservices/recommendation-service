# Building stage
FROM golang:1.19.3-alpine3.16 AS builder

WORKDIR /app
COPY . .

RUN go build -o main main.go


# Running stage
FROM alpine:3.16

WORKDIR /app

COPY --from=builder /app/main .

COPY app.env .
COPY start.sh .

RUN apk add --no-cache jq

RUN chmod +x ./start.sh

EXPOSE 8080

CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]