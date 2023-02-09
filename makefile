.PHONY: clean
clean:
	go clean -testcache

.PHONY: test
test:
	go test ./... -v

.PHONY: run
run:
	go run .