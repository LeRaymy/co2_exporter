FROM golang:1.18.2-alpine as builder

WORKDIR /go/co2_exporter

COPY . .

RUN go mod download
RUN go mod verify
RUN go build -o main ./cmd/co2_exporter

FROM golang:1.18.2-alpine

COPY --from=builder /go/co2_exporter/main /go/co2_exporter/main

EXPOSE 2112

CMD ["./co2_exporter/main"]