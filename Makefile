SHELL=/bin/bash

.PHONY: run
run:
	go run -race cmd/solocms/main.go run -config=./configs/config.dev.yaml

.PHONY: build
build:
	bash scripts/build.sh amd64

.PHONY: linter
linter:
	bash scripts/linter.sh

.PHONY: tests
tests:
	bash scripts/tests.sh

.PHONY: develop_up develop_down
develop_up:
	bash scripts/docker.sh docker_up
develop_down:
	bash scripts/docker.sh docker_down

.PHONY: ci
ci:
	bash scripts/ci.sh

deb:
	deb-builder build