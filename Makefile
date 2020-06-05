DOCKER_REGISTRY=radixdev.azurecr.io
VERSION=latest
IMAGE_NAME=$(DOCKER_REGISTRY)/radix-export-cost:$(VERSION)
DB_PASSWORD=a_password

# to deploy run:
# make deploy DB_PASSWORD=<sql_db_password>

build:
	docker build -t $(IMAGE_NAME) .

push:
	az acr login -n $(DOCKER_REGISTRY)
	docker push $(IMAGE_NAME)

deploy:
	helm upgrade --install radix-export-cost ./charts --set db.password=$(DB_PASSWORD)