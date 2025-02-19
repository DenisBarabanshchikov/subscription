run:
	godotenv -f .env go run ./cmd/main.go

doc:
	swag init -g cmd/main.go

test-unit:
	@echo "=== Begin Unit Tests ==="
	go test ./... --tags=unit -race -coverprofile coverage-unit.out -count=1
	@echo "=== End Unit Tests ==="

test-integration:
	@echo "=== Begin Integration Tests ==="
	godotenv -f .env go test ./... --tags=integration -race -coverprofile coverage-integration.out -count=1
	@echo "=== End Integration Tests ==="

check-test-build-flags:
	@FILES=$$(find . -name '*_test.go' -type f -exec awk 'NR == 1{ if ($$0 !~ /^\/\/go:build (unit|component|integration|scenario)$$/) print FILENAME; exit }' {} \;) ;\
	if [ -n "$$FILES" ]; then \
		echo 'found files with missing test build flags:'; \
		echo "$$FILES"; \
		exit 1; \
	fi

test: check-test-build-flags test-unit test-integration