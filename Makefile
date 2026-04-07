.PHONY: build clean sync test

# Sync skills/ to internal/templates/skills/ (for go:embed)
sync:
	@rm -rf internal/templates/skills
	@cp -r skills internal/templates/skills
	@cp FORGE.md internal/templates/FORGE.md
	@echo "✓ Synced skills and FORGE.md to internal/templates/"

# Build the CLI
build: sync
	@go build -o forge .
	@echo "✓ Built forge"

# Clean build artifacts
clean:
	@rm -f forge

# Run tests
test:
	@go test ./...

# Install to /usr/local/bin
install: build
	@sudo mv forge /usr/local/bin/
	@echo "✓ Installed forge to /usr/local/bin/"
