.PHONY: clean
clean:
	go clean -testcache

.PHONY: test
test:
	go test ./... -v -timeout 5s

.PHONY: test-non-verbose
test-non-verbose:
	go test ./... -timeout 5s

.PHONY: coverage-test
coverage-test:
	go test ./... -cover -timeout 5s

.PHONY: run
run:
	go run .