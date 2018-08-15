FROM golang:1.10.3-alpine3.8
ENV PORT 3000

WORKDIR /go/src/app
COPY . .
RUN apk add git
RUN go get github.com/nlopes/slack && \
    go get github.com/joho/godotenv && \
    go get gopkg.in/go-playground/webhooks.v5/github
EXPOSE 3000
CMD ["go","run","main.go"]
