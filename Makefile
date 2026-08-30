.PHONY: check format test test-race vet build tidy-check mod-verify module-graph

check: format tidy-check mod-verify module-graph test test-race vet

format:
	test -z "$$(find . -type f -name '*.go' -not -path './vendor/*' -print0 | xargs -0 gofmt -l)"

test:
	go test -mod=readonly -count=1 ./...

test-race:
	go test -mod=readonly -race -count=1 ./...

vet:
	go vet -mod=readonly ./...

build:
	go build -mod=readonly ./...

tidy-check:
	go mod tidy -diff

mod-verify:
	go mod verify

module-graph:
	test "$$(go list -mod=readonly -m all | wc -l | tr -d ' ')" -eq 1
