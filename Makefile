DOCKER_REGISTRY=radixdev.azurecr.io
VERSION=latest
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
IMAGE_NAME=$(DOCKER_REGISTRY)/radix-cost-allocation:$(BRANCH)-$(VERSION)

build:
	docker build -t $(IMAGE_NAME) .

push:
	az acr login -n $(DOCKER_REGISTRY)
	docker push $(IMAGE_NAME)

.PHONY: test
test:
	go test -cover `go list ./...`

.PHONY: lint
lint: bootstrap
	golangci-lint run --max-same-issues 0

.PHONY: mocks
mocks: bootstrap
	mockgen -source ./pkg/repository/repository.go -destination ./pkg/repository/mock/repository.go -package mock
	mockgen -source ./pkg/listers/limitrange.go -destination ./pkg/listers/mock/limitrange.go -package mock
	mockgen -source ./pkg/listers/node.go -destination ./pkg/listers/mock/node.go -package mock
	mockgen -source ./pkg/listers/pod.go -destination ./pkg/listers/mock/pod.go -package mock
	mockgen -source ./pkg/listers/radixregistration.go -destination ./pkg/listers/mock/radixregistration.go -package mock
	mockgen -source ./pkg/listers/containerbulkdto.go -destination ./pkg/listers/mock/containerbulkdto.go -package mock
	mockgen -source ./pkg/listers/nodebulkdto.go -destination ./pkg/listers/mock/nodebulkdto.go -package mock

HAS_GOLANGCI_LINT := $(shell command -v golangci-lint;)
HAS_MOCKGEN       := $(shell command -v mockgen;)

bootstrap:
ifndef HAS_GOLANGCI_LINT
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2
endif
ifndef HAS_MOCKGEN
	go install github.com/golang/mock/mockgen@v1.6.0
endif
