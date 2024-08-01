.PHONY: run build race tidy test

run:
	go run ./cmd/goker/main.go

build:
	go build ./cmd/goker/main.go

race:
	go run --race cmd/goker/main.go

tidy:
	go mod tidy

test:
	go clean -testcache
	go test -v ./... | sed "/PASS/s//$$(printf "\033[32mPASS\033[0m")/" | sed "/FAIL/s//$$(printf "\033[31mFAIL\033[0m")/"
