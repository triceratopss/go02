FROM golang:1.22-bookworm AS debug

WORKDIR /app

COPY ./go.* /app/
RUN go mod download
RUN go install github.com/air-verse/air@latest

COPY . /app

CMD ["air", "-c", ".air.toml"]

# build continer
FROM golang:1.22-bookworm AS builder

WORKDIR /tmp/app

COPY go.* ./
RUN go mod download

COPY . ./
RUN go build -o "go02" -ldflags="-w -s"


# runtime continer
FROM gcr.io/distroless/base-debian12

COPY --from=builder /tmp/app/go02 /root/go02

CMD ["/root/go02"]
