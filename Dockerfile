FROM golang:1.14-alpine as builder

RUN apk update && \
    apk add ca-certificates curl git  && \
    apk add --no-cache gcc musl-dev && \
    go get -u golang.org/x/lint/golint github.com/frapposelli/wwhrd

WORKDIR /go/src/github.com/equinor/radix-export-cost/

# Install project dependencies
COPY go.mod go.sum ./
RUN go mod download

# Check dependency licenses using https://github.com/frapposelli/wwhrd
COPY .wwhrd.yml ./
RUN wwhrd -q check

# Copy project code
COPY . .

# run tests and linting
RUN golint `go list ./...` && \
    go vet `go list ./...` && \
    CGO_ENABLED=0 GOOS=linux go test `go list ./...`

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -a -installsuffix cgo -o ./rootfs/radix-export-cost
RUN adduser -D -g '' radix-export-cost


# Run operator
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /go/src/github.com/equinor/radix-export-cost/rootfs/radix-export-cost /usr/local/bin/radix-export-cost
USER radix-export-cost

ENTRYPOINT ["/usr/local/bin/radix-export-cost"]