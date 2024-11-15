FROM golang:1.23-alpine AS builder
RUN apk add build-base
WORKDIR /app
COPY . /app
RUN go build -o stori main.go

FROM builder AS runner
ENTRYPOINT ["./stori"]
