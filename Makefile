DOCKER_REGISTRY=radixdev.azurecr.io
VERSION=latest
IMAGE_NAME=$(DOCKER_REGISTRY)/radix-cost-allocation:$(VERSION)
DB_PASSWORD=a_password

# to deploy run: "make deploy DB_PASSWORD=<sql_db_password>"

# to deploy db: "make deploy-azure DB_PASSWORD=<sql_db_password>"

build:
	docker build -t $(IMAGE_NAME) .

push:
	az acr login -n $(DOCKER_REGISTRY)
	docker push $(IMAGE_NAME)

deploy:
	helm upgrade --install radix-cost-allocation ./charts --set db.password=$(DB_PASSWORD)

deploy-azure:
	az deployment group create --resource-group common --template-file ./azure-infrastructure/azuredeploy.json --parameters sqlAdministratorLoginPassword=$(DB_PASSWORD)

.PHONY: test
test:
	go test -cover `go list ./...`

.PHONY: lint
lint: bootstrap
	golangci-lint run --max-same-issues 0

HAS_GOLANGCI_LINT := $(shell command -v golangci-lint;)

bootstrap:
ifndef HAS_GOLANGCI_LINT
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.55.2
endif
