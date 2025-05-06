FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app/main.go

FROM alpine:3.21

WORKDIR /app


COPY --from=builder /app/app .
COPY --from=builder /app/env.yaml .
COPY --from=builder /app/templates/ templates/

RUN apk add --no-cache tzdata
RUN chmod +x ./app && apk --no-cache add dumb-init

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["./app"]