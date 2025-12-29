# ── project metadata ────────────────────────────────────────────────────────────
name        := gitr
pkg         := github.com/swarupdonepudi/gitr
build_dir   := dist
LDFLAGS     := -ldflags "-X $(pkg)/cmd/gitr/root.VersionLabel=$$(git describe --tags --always --dirty)"

# ── helper vars ────────────────────────────────────────────────────────────────
build_cmd   := go build $(LDFLAGS)

# ── quality / housekeeping ─────────────────────────────────────────────────────
.PHONY: deps vet fmt test clean
deps:          ## download & tidy modules
	go mod download
	go mod tidy

vet:           ## go vet
	go vet ./...

fmt:           ## go fmt
	go fmt ./...

test: vet      ## run tests with race detector
	go test -race -v -count=1 ./...

clean:         ## remove build artifacts
	rm -rf $(build_dir)

# ── build ─────────────────────────────────────────────────────────────────────
.PHONY: build build-cli build-site
build: build-cli build-site ## build CLI and website

build-cli: deps fmt vet ## build the Go CLI binary
	$(build_cmd) -o $(build_dir)/$(name) .

build-site: ## build the website
	cd site && NODE_NO_WARNINGS=1 yarn install
	cd site && NODE_NO_WARNINGS=1 yarn build

# ── local utility ──────────────────────────────────────────────────────────────
.PHONY: snapshot local
snapshot: deps ## build a local snapshot using GoReleaser
	goreleaser release --snapshot --clean --skip=publish

local: deps fmt vet ## build and install binary to ~/bin
	$(build_cmd) -o $(build_dir)/$(name) .
	install -m 0755 $(build_dir)/$(name) $(HOME)/bin/$(name)

# ── release tagging ────────────────────────────────────────────────────────────
.PHONY: release build-check next-version
build-check:   ## quick compile to verify build
	go build -o /dev/null .

# bump: major, minor, or patch (default)
bump ?= patch

next-version:  ## show what the next version would be
	@latest=$$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0"); \
	major=$$(echo $$latest | sed 's/v//' | cut -d. -f1); \
	minor=$$(echo $$latest | sed 's/v//' | cut -d. -f2); \
	patch=$$(echo $$latest | sed 's/v//' | cut -d. -f3); \
	case "$(bump)" in \
		major) major=$$((major + 1)); minor=0; patch=0 ;; \
		minor) minor=$$((minor + 1)); patch=0 ;; \
		patch) patch=$$((patch + 1)) ;; \
		*) echo "Invalid bump type: $(bump). Use major, minor, or patch"; exit 1 ;; \
	esac; \
	echo "v$$major.$$minor.$$patch"

release: test build-check ## auto-bump version, tag & push (bump=major|minor|patch, default: patch)
	@latest=$$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0"); \
	major=$$(echo $$latest | sed 's/v//' | cut -d. -f1); \
	minor=$$(echo $$latest | sed 's/v//' | cut -d. -f2); \
	patch=$$(echo $$latest | sed 's/v//' | cut -d. -f3); \
	case "$(bump)" in \
		major) major=$$((major + 1)); minor=0; patch=0 ;; \
		minor) minor=$$((minor + 1)); patch=0 ;; \
		patch) patch=$$((patch + 1)) ;; \
		*) echo "Invalid bump type: $(bump). Use major, minor, or patch"; exit 1 ;; \
	esac; \
	version="v$$major.$$minor.$$patch"; \
	echo "Current version: $$latest"; \
	echo "Releasing: $$version ($(bump) bump)"; \
	git tag -a $$version -m "$$version"; \
	git push origin $$version

# ── default target ─────────────────────────────────────────────────────────────
.DEFAULT_GOAL := test


.PHONY: develop-site
develop-site:
	cd site && NODE_NO_WARNINGS=1 yarn install
	cd site && NODE_NO_WARNINGS=1 yarn dev

.PHONY: preview-site
preview-site:
	cd site && NODE_NO_WARNINGS=1 yarn install
	cd site && NODE_NO_WARNINGS=1 yarn build:serve
