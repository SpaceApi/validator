FROM golang:1.15-alpine as builder
RUN apk --no-cache add git build-base
WORKDIR /app
COPY . .
RUN go get -d  ./...
RUN go generate
RUN go test ./...
RUN go install  ./...


FROM alpine:latest
RUN apk --no-cache add ca-certificates curl
WORKDIR /app
COPY --from=builder /go/bin/validator /usr/local/bin/validator
EXPOSE 8080

HEALTHCHECK --start-period=5s CMD curl --fail http://localhost:8080/v1/ || exit 1

RUN adduser app -S -u 142
USER app

CMD ["validator"]
