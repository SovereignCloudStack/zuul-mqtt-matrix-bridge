FROM --platform=${BUILDPLATFORM} golang:1.21 as builder

ARG TARGETOS
ARG TARGETARCH

WORKDIR /build
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
COPY pkg/ pkg/
COPY templates/ templates/
COPY main.go main.go
RUN ls -lt 
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -o zuul-mqtt-matrix-bridge main.go

FROM --platform=${TARGETPLATFORM} scratch
WORKDIR /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/zuul-mqtt-matrix-bridge .
COPY --from=builder /build/templates /templates

USER 65532:65532

ENTRYPOINT ["/zuul-mqtt-matrix-bridge"]
