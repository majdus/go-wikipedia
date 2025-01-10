.DEFAULT_GOAL := precommit
IMPI=impi
REPO=github.com/majdus/go-wikipedia
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)

ALL_GO_MOD_DIRS := $(shell find . -type f -name 'go.mod' -exec dirname {} \; | sort)

.PHONY: precommit
precommit: install-tools build lint impi fmt test

.PHONY: impi
impi:
	@$(IMPI) --local $(REPO) --scheme stdThirdPartyLocal ./...


.PHONY: lint
lint:
	set -e; for dir in $(ALL_GO_MOD_DIRS); do \
	  echo "golangci-lint in $${dir}"; \
	  (cd "$${dir}" && \
	    golangci-lint run --fix && \
	    golangci-lint run); \
	done
	set -e; for dir in $(ALL_GO_MOD_DIRS); do \
	  echo "go mod tidy in $${dir}"; \
	  (cd "$${dir}" && \
	    go mod tidy); \
	done

.PHONY: build
build:
	@go build -v ./...
	@echo "build success"

.PHONY: fmt
fmt:
	gofmt -w -s .
	goimports -w -local $(REPO) ./

.PHONY: test
test: unit-test

.PHONY: unit-test
unit-test:
	go test -v ./...

.PHONY: install-tools
install-tools:
	@which golangci-lint > /dev/null || (echo 'install golangci-lint' && go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3)
	@which misspell > /dev/null || (echo 'install misspell' && go install github.com/client9/misspell/cmd/misspell@latest) # 检查单词拼写错误
	@which addlicense > /dev/null || (echo 'install google-addlicense' && go install github.com/google/addlicense@latest) # 检查或添加license到每个文件
	@which impi > /dev/null || (echo 'install impi' && go install github.com/pavius/impi/cmd/impi@latest) # 检查go语言import分组是否符合规范
	@which goimports > /dev/null || (echo 'install goimports' && go install golang.org/x/tools/cmd/goimports@latest)

