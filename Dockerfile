FROM golang:1.23-alpine AS builder
WORKDIR /src
COPY api api
COPY internal internal
COPY pkg pkg
COPY main.go main.go
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download && go build -o watermark main.go

FROM alpine:3.17
WORKDIR /app
COPY --from=builder /src/watermark /app/watermark
EXPOSE 8081 8082
ENTRYPOINT ["/app/watermark"]