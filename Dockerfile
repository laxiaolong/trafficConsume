FROM golang:alpine AS builder

WORKDIR /build
COPY . .

RUN go mod tidy
RUN go build -trimpath -v -ldflags "-X main.date=$(date -Iseconds)"

FROM alpine

RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /build/trafficConsume /app/trafficConsume
ENTRYPOINT ["/app/trafficConsume"]