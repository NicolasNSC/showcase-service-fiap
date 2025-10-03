all: build run

build: 
	go build -o build/bin/showcase-service-fiap ./cmd/showcase-service-fiap 

run:
	./build/bin/showcase-service-fiap

test: 
	go test -covermode=atomic -coverprofile=coverage.out `go list ./... | grep -v mocks | grep -v cmd | grep -v testdata`

cov: test
	go tool cover -html=coverage.out

gen: 
	go generate ./...

swagger:
	swag init -g cmd/showcase-service-fiap/main.go -o ./docs --parseDependency --parseInternal

docker-up:
	docker-compose up -d --build

docker-down:
	docker-compose down -v