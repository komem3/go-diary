

init: ## init
	go mod download

test: ## test
	go test ./...

test/intergration: ## intergration test
	cd ./testdata/intergration && \
	./intergration_test.sh

gen: ## generate task
	statik -src=./template

install: ## install diary
	cd cmd/diary && go install .

clean: ## clean
	go clean && go mod tidy

build: build-linux	build-mac	build-win  ## do all build task

build-linux: ## build linux 64bit binary
	GOOS=linux GOARCH=amd64 go build -o build/linux-amd64/diary ./cmd/diary/main.go
build-mac: ## build mac os 64bit binary
	GOOS=darwin GOARCH=amd64 go build -o build/darwin-amd64/diary ./cmd/diary//main.go
build-win: ## build windows os 64bit binary
	GOOS=windows GOARCH=amd64 go build -o build/windows-amd64/diary ./cmd/diary/main.go
