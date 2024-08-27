include build.mk

export GO111MODULE=on
export GOSUMDB=off

BUILD_ENV_PARAMS:=CGO_ENABLED=0

space := $(subst ,, )
CURDIR_ESCAPE:=$(subst $(space),\ ,$(CURDIR))

LOCAL_BIN:=$(CURDIR_ESCAPE)/bin
LINT_BIN:=$(LOCAL_BIN)/golangci-lint
LINT_VERSION:=1.60.3

###### TEST ######
.PHONY: test
test:
	go test ./... -count=1 -timeout=60s -v -short
###### TEST ######

###### LINT ######
.PHONY: install-lint
install-lint:
ifeq ($(wildcard $(LINT_BIN)),)
	$(info Installing golangci-lint v$(LINT_VERSION))
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v$(LINT_VERSION)
# Устанавливаем текущий путь для исполняемого файла линтера.
else
	$(info Golangci-lint is already installed to $(LINT_VERSION))
endif


PHONY: lint
lint: install-lint
	$(info Running lint against changed files...)
	$(LINT_BIN) run \
		--new-from-rev=origin/master \
		--config=.golangci.yml \
		./...

PHONY: lint-full
lint-full: install-lint
	$(info Running lint against all project files...)
	$(LINT_BIN) run \
		--config=.golangci.yml \
		./...
###### LINT ######