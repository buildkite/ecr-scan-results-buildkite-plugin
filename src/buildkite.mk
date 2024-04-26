all: clean lint test ecrscanresults

ecrscanresults:
	CGO_ENABLED=0 go build -o ecrscanresults -trimpath -mod=readonly -ldflags="-s -w -X main.version=$(shell git describe --always)" .

.PHONY: clean
clean:
	@rm -f ecrscanresults
	@rm -rf dist

.PHONY: lint
lint:
	go vet ./...

.PHONY: test
test:
	go test ./...

cover.out:
	go test ./... -coverprofile=cover.out

.PHONY: coverage
coverage: cover.out
	go tool cover -html=cover.out

.PHONY: tidy
tidy:
	go mod tidy