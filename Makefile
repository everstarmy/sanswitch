.PHONY: check format test test-shuffle test-race coverage vet build tidy-check mod-verify module-graph

check: format tidy-check mod-verify module-graph test test-shuffle test-race coverage vet

format:
	test -z "$$(find . -type f -name '*.go' -not -path './vendor/*' -print0 | xargs -0 gofmt -l)"

test:
	go test -mod=readonly -count=1 ./...

test-shuffle:
	go test -mod=readonly -shuffle=on -count=1 ./...

test-race:
	go test -mod=readonly -race -count=1 ./...

coverage:
	profile="$$(mktemp)"; trap 'rm -f "$$profile"' EXIT; \
		go test -mod=readonly -coverprofile="$$profile" ./...; \
		go tool cover -func="$$profile"; \
	go tool cover -func="$$profile" | awk '/^total:/ { gsub("%", "", $$3); if ($$3 + 0 < 75) exit 1 }'

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
