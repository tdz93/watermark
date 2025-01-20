FROM golang:1.23-alpine AS builder
WORKDIR /src
COPY api/v1/pb/watermark api/v1/pb/watermark
COPY cmd/watermark cmd/watermark
COPY internal internal
COPY pkg/watermark pkg/watermark
COPY go.mod go.mod
COPY go.sum go.sum
RUN go build -o watermark cmd/watermark/watermark.go

FROM alpine:3.17
WORKDIR /app
COPY --from=builder /src/watermark /app/watermark
EXPOSE 8081 8082
ENTRYPOINT ["/app/watermark"]