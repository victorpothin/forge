.PHONY: build clean sync test install release

# Sync skills/ to internal/templates/skills/ (for go:embed)
sync:
	@rm -rf internal/templates/skills
	@cp -r skills internal/templates/skills
	@cp FORGE.md internal/templates/FORGE.md
	@echo "✓ Synced skills and FORGE.md to internal/templates/"

# Build with version from git tag
VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null | sed 's/^v//')
ifeq ($(VERSION),)
VERSION := dev
endif

LDFLAGS := -ldflags "-s -w -X github.com/victorpothin/forge/cmd.version=$(VERSION)"

build: sync
	@go build $(LDFLAGS) -o forge .
	@echo "✓ Built forge $(VERSION)"

# Clean build artifacts
clean:
	@rm -f forge

# Run tests
test:
	@go test ./...

# Install to /usr/local/bin
install: build
	@sudo mv forge /usr/local/bin/
	@echo "✓ Installed forge $(VERSION) to /usr/local/bin/"

# Release — clean build with current git tag
release: sync
	@git fetch --tags 2>/dev/null || true
	@LATEST_TAG=$$(git describe --tags --abbrev=0 2>/dev/null); \
	if [ -z "$$LATEST_TAG" ]; then echo "No git tag found"; exit 1; fi; \
	VERSION=$$(echo $$LATEST_TAG | sed 's/^v//'); \
	echo "Building $$LATEST_TAG..."; \
	go build -ldflags "-s -w -X github.com/victorpothin/forge/cmd.version=$$VERSION" -o forge .; \
	./forge --version
