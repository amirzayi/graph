FROM golang:1.27.1-alpine3.24 as builder
LABEL authors="amirzayi"
WORKDIR /app
COPY . .
RUN go mod download && CGO_ENABLED=0 go build -o bin -ldflags '-s -w' .

FROM alpine
COPY --from=builder /app/bin .
EXPOSE 8010
ENTRYPOINT ["./bin"]
