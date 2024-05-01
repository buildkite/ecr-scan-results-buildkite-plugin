all: clean lint test src/ecrscanresults

src/ecrscanresults:
	CGO_ENABLED=0 go -C src build -o ecrscanresults -trimpath -mod=readonly -ldflags="-s -w -X main.version=$(shell git describe --always)" .

.PHONY: clean
clean:
	@rm -f src/ecrscanresults
	@rm -rf src/dist

.PHONY: lint
lint: tidy
	# go vet ./...
	golangci-lint run --out-format=github-actions --path-prefix=src -v --timeout=2m

.PHONY: test
test:
	go -C src test ./...

cover.out:
	go -C src test ./... -coverprofile=cover.out

.PHONY: coverage
coverage: cover.out
	go -C src tool cover -html=cover.out

.PHONY: tidy
tidy:
	go -C src mod tidy