FROM --platform=$BUILDPLATFORM docker.io/golang:1.25.5-alpine3.23 AS builder

ARG TARGETARCH

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=$TARGETARCH

WORKDIR /src

# Install project dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy project code
COPY . .

RUN go build -ldflags="-s -w" -o /build/radix-cost-allocation

# Final stage, ref https://github.com/GoogleContainerTools/distroless/blob/main/base/README.md for distroless
FROM gcr.io/distroless/static
WORKDIR /app
COPY --from=builder /build/radix-cost-allocation .
USER 1000
ENTRYPOINT ["/app/radix-cost-allocation"]
