.PHONY: build run test clean docker

build:
	go build -o bin/snmp-zte ./cmd/api

run:
	go run ./cmd/api

test:
	go test -v ./...

clean:
	rm -rf bin/

docker:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f

deps:
	go mod download
	go mod tidy
