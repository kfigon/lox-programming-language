.PHONY: clean
clean:
	go clean -testcache

.PHONY: test
test:
	go test ./... -v

.PHONY: coverage-test
coverage-test:
	go test ./... -cover

.PHONY: run
run:
	go run .