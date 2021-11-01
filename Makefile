LOCAL_BIN ?=./bin

version=1.0.0
container_name=PickupStats

.PHONY: build
build:
	docker build -t condensedtea/pickupstats:latest -t condensedtea/pickupstats:$(version) .


.PHONY: app
app:
	CGO_ENABLED=0 go build -o "$(LOCAL_BIN)/app" ./app

PHONY: run
run:
	docker run --network=host --rm -d --name=$(container_name) condensedtea/pickupstats:latest

PHONY: down
down:
	docker kill $(container_name)