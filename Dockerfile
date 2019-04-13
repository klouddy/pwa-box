FROM golang:alpine AS builder
RUN apk update && apk add --no-cache git
RUN adduser -D -g '' appuser
RUN mkdir /src
WORKDIR /src
ADD *.go ./
ADD resources/ resources/
ADD go.* ./
RUN go mod download
RUN go build -o /go/bin/pwa-box

FROM alpine:3.7
EXPOSE 9898
ENV PATH=/go/bin:$PATH
RUN mkdir /pwa-box
WORKDIR /pwa-box

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /src/resources/default-config.json config.json
COPY --from=builder /src/resources/default-static-site/ static/
COPY --from=builder /go/bin/pwa-box /go/bin/pwa-box
USER appuser
CMD ["pwa-box"]
