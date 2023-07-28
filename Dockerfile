FROM golang:1.20-alpine3.18 as builder

RUN apk update && \
    apk add ca-certificates curl git && \
    apk add --no-cache gcc musl-dev
RUN go install honnef.co/go/tools/cmd/staticcheck@2023.1.3

WORKDIR /go/src/github.com/equinor/radix-cost-allocation/

# Install project dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy project code
COPY . .

# run tests and linting
RUN staticcheck ./... && \
    go vet ./... && \
    CGO_ENABLED=0 GOOS=linux go test `go list ./...`

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -a -installsuffix cgo -o ./rootfs/radix-cost-allocation
RUN addgroup -S -g 1000 radix-cost-allocation
RUN adduser -S -u 1000 -G radix-cost-allocation radix-cost-allocation

# Run operator
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/src/github.com/equinor/radix-cost-allocation/rootfs/radix-cost-allocation /usr/local/bin/radix-cost-allocation
USER 1000

ENTRYPOINT ["/usr/local/bin/radix-cost-allocation"]
