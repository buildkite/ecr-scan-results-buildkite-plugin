all: clean lint test src/ecrscanresults

src/ecrscanresults:
	CGO_ENABLED=0 go -C src build -o ecrscanresults -trimpath -mod=readonly -ldflags="-s -w -X main.version=$(shell git describe --always)" .

.PHONY: clean
clean:
	@rm -f src/ecrscanresults
	@rm -rf src/dist

.PHONY: lint
lint:
	go -C src mod tidy
	(cd src && golangci-lint run --out-format=github-actions --path-prefix=src -v --timeout=2m)

cover.out:
	go -C src test ./... -coverprofile=cover.out

.PHONY: coverage
coverage: cover.out
	go -C src tool cover -html=cover.out

# Delegate these targets to the upstream makefile, src/Makefile
.PHONY: mod test test-ci ensure-deps
mod test test-ci ensure-deps:
	$(MAKE) -C src $@