LOCAL_BIN ?=./bin

version=1.0.0
container_name=PickupStats

LDFLAGS = "-X 'main.Version=$(version)'"

.PHONY: build
build:
	docker build -t condensedtea/pickupstats:latest -t condensedtea/pickupstats:$(version) .


.PHONY: app
app:
	CGO_ENABLED=0 go build -ldflags=$(LDFLAGS) -o "$(LOCAL_BIN)/app" ./app

PHONY: run
run:
	docker run --network=host --rm -d --name=$(container_name) condensedtea/pickupstats:latest

PHONY: down
down:
	docker kill $(container_name)

.PHONY: docs
docs:
	swag init -g main.go -d app/,pkg/api,pkg/db --output docs/