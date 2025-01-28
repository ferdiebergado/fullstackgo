.PHONY: $(wildcard *)

%:
	@true

test:
	go test -v ./...
