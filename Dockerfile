FROM golang:1.21 as builder

WORKDIR /build
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY pkg/ pkg/
COPY templates/ templates/
COPY main.go main.go
RUN ls -lt 
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o zuul-mqtt-matrix-bridge main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /build/zuul-mqtt-matrix-bridge .
COPY --from=builder /build/templates /templates

USER 65532:65532

ENTRYPOINT ["/zuul-mqtt-matrix-bridge"]
