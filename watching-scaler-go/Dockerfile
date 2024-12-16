FROM golang:1.23 as builder
ENV CGO_ENABLED=0

WORKDIR /build
COPY . .

RUN go mod download
RUN go build -o app

FROM alpine:latest as runtime
COPY --from=builder /build/app /usr/local/bin/app

CMD ["app"]