FROM golang:1.19-alpine

ENV GOPATH=/

COPY ./ ./

RUN go mod download
RUN go build -o balance-service ./cmd/main.go

CMD ["./balance-service"]