PWD=$(shell pwd)
export GOPATH=${HOME}

ifndef TARGET
	TARGET=$(shell go list ./... | grep -v /vendor/)
endif

.PHONY: create-migration get build start run test-short test test-cli lint docs

get:
	go get -v ${TARGET}
	go get -u -d github.com/mattes/migrate/cli github.com/lib/pq
	go get -u github.com/kardianos/govendor
	go build -tags 'postgres' -o ./bin/migrate github.com/mattes/migrate/cli

BUILD_ARGS=-v
build-internal:
	go get -t ${BUILD_ARGS} ${TARGET}

start-internal:
	~/bin/cryptocades-backend

test-short-internal:
	go test -short ${TARGET}

test-internal:
	go test ${TARGET}

test-cover-internal:
	./scripts/cover.sh

vet-internal:
	go vet ${TARGET}

fmt-check-internal:
	go fmt ${TARGET}

docs:
	$(MAKE) -C docs build-html build-swagger

codecov:
	curl -s https://codecov.io/bash > /tmp/codecov.bash
	bash /tmp/codecov.bash -t 154e8023-2188-4fde-b714-99967c6ead80 -f coverage.txt -X fix

wait-for-postgres:
	@./scripts/wait-for-postgres.sh

create-migration:
	@./bin/migrate create -ext sql -dir migrations ${name}
	#@docker-compose run --rm --no-deps cryptocades-backend ./bin/migrate create -ext sql -dir migrations ${name}

sh:
	@docker-compose run --rm --no-deps cryptocades-backend /bin/sh

build:
	@docker-compose run --rm --no-deps cryptocades-backend make build-internal

# this leave postgres running
start:
	@docker-compose run --rm -p 8080:8080 cryptocades-backend make wait-for-postgres start-internal

run: build start

test-short:
	@docker-compose run --rm --no-deps -e TARGET=$$TARGET cryptocades-backend make test-short-internal

test:
	@./scripts/run.sh cryptocades-backend make wait-for-postgres test-internal

test-cover:
	@./scripts/run.sh cryptocades-backend make wait-for-postgres test-cover-internal

lint-internal:
	make fmt-check-internal vet-internal

lint:
	@docker-compose run --rm -e TARGET="${TARGET}" --no-deps cryptocades-backend make lint-internal 
