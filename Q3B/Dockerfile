FROM golang:1.10.3-alpine3.8
ENV PORT 8000
EXPOSE 8000

WORKDIR /go/src/app
COPY . .
RUN apk add git
RUN go get github.com/gorilla/mux && \
    go get github.com/joho/godotenv 
CMD ["go","run","main.go"]
