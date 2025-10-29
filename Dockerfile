FROM golang:1.23 AS builder
WORKDIR /usr/src/app

COPY . .

RUN go mod download
RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o subscriptions main.go

FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache bash

COPY --from=builder /usr/src/app/subscriptions .
COPY --from=builder /usr/src/app/wait-for-it.sh .
COPY --from=builder /usr/src/app/GO/database/migrations ./GO/database/migrations
COPY --from=builder /usr/src/app/.env ./

RUN chmod +x wait-for-it.sh

EXPOSE 8080
ENV APP_PORT=8080

CMD ["bash", "./wait-for-it.sh", "db:5432", "--", "./subscriptions"]