.PHONY: dockerbuild
dockerbuild:
	docker build -t codacy-usage-report .

.PHONY: build
build:
	./scripts/cross-build.sh

.PHONY: run
run:
	go run main.go

.PHONY: test
test:
	go test ./...
