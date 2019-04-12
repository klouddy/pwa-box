FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
RUN adduser -D -g '' appuser
RUN mkdir /src
WORKDIR /src
ADD *.go ./
ADD config.json config.json
ADD go.* ./
RUN go mod download
RUN go build -o /go/bin/pwa-box

FROM alpine:3.7

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /src/config.json config.json
COPY --from=builder /go/bin/pwa-box /go/bin/pwa-box
USER appuser
CMD ["/go/bin/pwa-box"]
