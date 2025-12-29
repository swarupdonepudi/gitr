pkg=github.com/swarupdonepudi/gitr
LDFLAGS=-ldflags "-X ${pkg}/pkg/version.Version=${v}"
build_cmd=go build -v ${LDFLAGS}

.PHONY: deps
deps:
	go mod download
.PHONY: vet
vet:
	go vet ./...

.PHONY: fmt
fmt:
	go fmt ./...
.PHONY: build
build: deps vet fmt
	env GOOS=darwin ${build_cmd} -o bin/gitr-darwin main.go
	env GOOS=darwin GOARCH=amd64 ${build_cmd} -o bin/gitr-darwin-amd64 main.go
	env GOOS=darwin GOARCH=arm64 ${build_cmd} -o bin/gitr-darwin-arm64 main.go
.PHONY: checksum

checksum: build
	openssl dgst -sha256 bin/gitr-darwin-amd64
	openssl dgst -sha256 bin/gitr-darwin-arm64

.PHONY: setup-tests
setup-tests:
	mv pkg/git/git_test_data/r1-no-remote/.git-temp pkg/git/git_test_data/r1-no-remote/.git
	mv pkg/git/git_test_data/r2-with-remote/.git-temp pkg/git/git_test_data/r2-with-remote/.git
	mv pkg/git/git_test_data/r3-with-remote-custom-branch/.git-temp pkg/git/git_test_data/r3-with-remote-custom-branch/.git
.PHONY: execute-tests
execute-tests:
	go test -v -coverpkg github.com/swarupdonepudi/gitr/internal/..  -cover ./... -coverprofile=internal.cov || true
	go test -v -coverpkg github.com/swarupdonepudi/gitr/pkg/...  -cover ./... -coverprofile=pkg.cov || true
.PHONY: cleanup-tests
cleanup-tests:
	mv pkg/git/git_test_data/r1-no-remote/.git pkg/git/git_test_data/r1-no-remote/.git-temp
	mv pkg/git/git_test_data/r2-with-remote/.git pkg/git/git_test_data/r2-with-remote/.git-temp
	mv pkg/git/git_test_data/r3-with-remote-custom-branch/.git pkg/git/git_test_data/r3-with-remote-custom-branch/.git-temp
.PHONY: test
test: setup-tests execute-tests cleanup-tests
.PHONY: analyze-tests
analyze-tests:
	go tool cover -func=internal.cov
	go tool cover -func=pkg.cov
.PHONY: local
local: build
	sudo cp bin/gitr-darwin /usr/local/bin/gitr
