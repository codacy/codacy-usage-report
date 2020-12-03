.PHONY: dockerbuild
dockerbuild:
	docker build -t codacy-usage-report .

.PHONY: build
build:
	go build -o bin/codacy-usage-report .

.PHONY: run
run:
	go run main.go

.PHONY: test
test:
	go test ./...
