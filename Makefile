help: ## help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

init: ## init
	go mod download && \
	go get golang.org/x/tools/cmd/goimports && \
	rm -r node_modules && \
	npm ci && \
	cp ./hack/precommit.sh .git/hooks/pre-commit

# Test
test: test_unit test_intergration ## exec unit test and intergration test
test_unit: ## exec unit test
	go test ./...
test_intergration: ## exec intergration test
	make install && \
	cd ./testdata/intergration && \
	./intergration_test.sh

gen: ## generate files
	statik -src=./template

install: ## install diary
	cd cmd/diary && go install .

clean: ## clean
	go clean && go mod tidy

# Static analysis
lint: lint_go lint_md ## lint source code
lint_go: ## go lint
	golangci-lint run ./... --disable-all \
	-E govet -E errcheck -E staticcheck -E unused -E gosimple \
	-E structcheck -E varcheck -E ineffassign -E deadcode -E typecheck \
	-E golint -E interfacer -E unconvert -E dupl -E goconst \
	-E asciicheck -E gofmt -E goimports -E misspell -E lll -E unparam \
	-E prealloc -E gocritic -E gochecknoinits -E whitespace -E gomnd \
	-E goerr113 -E gomodguard -E godot -E testpackage
lint_md: ## markdown lint
	npm run lint

# Build
build: build-linux	build-mac	build-win  ## do all build task
build-linux: ## build linux 64bit binary
	GOOS=linux GOARCH=amd64 go build -o build/linux-amd64/diary ./cmd/diary/main.go
build-mac: ## build mac os 64bit binary
	GOOS=darwin GOARCH=amd64 go build -o build/darwin-amd64/diary ./cmd/diary//main.go
build-win: ## build windows os 64bit binary
	GOOS=windows GOARCH=amd64 go build -o build/windows-amd64/diary ./cmd/diary/main.go
