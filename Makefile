.PHONY: $(wildcard *)

%:
	@true

test:
	go test -v ./...

gen:
	@command -v mockgen >/dev/null || go install go.uber.org/mock/mockgen@latest
	go generate -v ./...
