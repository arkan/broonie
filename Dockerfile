FROM golang:1.15
WORKDIR /go/src/github.com/arkan/telegram_memories_bot
RUN go get -d -v github.com/go-telegram-bot-api/telegram-bot-api
COPY main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=0 /go/src/github.com/arkan/telegram_memories_bot/app /app/bot
RUN mkdir /app/telegram && chmod -R 777 /app
CMD ["./bot"]