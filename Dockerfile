FROM golang:1.25.0-alpine as builder

WORKDIR /app
COPY . /app
RUN go build -o /app/main cmd/app/main.go

FROM alpine:3
WORKDIR /app
COPY --from=builder /app/main /app
COPY db /app/db
RUN chmod +x main
RUN touch .env
CMD ./main
