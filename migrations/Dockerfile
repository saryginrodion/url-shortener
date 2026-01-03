FROM golang:1.25.5-alpine AS builder
ENV GOPATH=/go
ENV GOBIN=/go/bin

RUN go install github.com/pressly/goose/v3/cmd/goose@v3.26.0

FROM scratch

COPY --from=builder /go/bin/goose /goose
COPY migrations/ /migrations/

WORKDIR /migrations

CMD ["/goose", "status"]
