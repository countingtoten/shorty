FROM golang:1.10-alpine as builder

WORKDIR /go/src/github.com/countingtoten/shorty/

COPY . .

RUN apk --update --no-cache add git

RUN go get -u github.com/golang/dep/cmd/dep

RUN dep ensure

RUN go test -cover ./...

RUN CGO_ENABLED=0 GOOS=linux go build -o shorty ./cmd/shorty/

FROM alpine:latest

RUN apk --update --no-cache add ca-certificates

WORKDIR /app/

COPY --from=builder /go/src/github.com/countingtoten/shorty/shorty .

CMD ["./shorty"]
