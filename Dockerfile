FROM golang:1.13-alpine as builder
RUN apk --no-cache add git build-base
WORKDIR /app
COPY . .
RUN go get -d  ./...
RUN go generate
RUN go test ./...
RUN go install  ./...


FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /go/bin/validator /usr/local/bin/validator
EXPOSE 8080

RUN adduser app -S -u 142
USER app

CMD ["validator"]
