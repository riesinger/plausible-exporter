FROM golang:1.19-alpine as builder

WORKDIR /go/src/app

# Install dependencies first for better caching
COPY go.mod go.sum /go/src/app/
RUN go get -v ./...

# Do a completely static build
COPY . .
RUN CGO_ENABLED=0 go install -ldflags '-s -w -extldflags "-static"' -tags timetzdata ./cmd
RUN ls -l /go/bin

FROM scratch as runner

# Copy the built binary
COPY --from=builder /go/bin/cmd /plausible-exporter
# Copy the CA root certificats from the latest alpine image
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT [ "/plausible-exporter" ]
