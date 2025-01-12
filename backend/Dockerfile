FROM golang:1.23-alpine

WORKDIR /go/src/app

COPY . .
ENV CGO_ENABLED=1
RUN apk add --no-cache gcc musl-dev
RUN go mod download
RUN go build -ldflags='-s -w -extldflags "-static"' -o main main.go

CMD ["./main"]