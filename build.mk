GO_VERSION_SHORT:=$(shell echo `go version` | sed -E 's/.* go(.*) .*/\1/g')
ifneq ("1.23","$(shell printf "$(GO_VERSION_SHORT)\n1.23" | sort -V | head -1)")
	$(error NEED GO VERSION >= 1.23. Found: $(GO_VERSION_SHORT))
endif

##################### PROJECT RELATED VARIABLES #####################
ifndef GIT_BRANCH
	GIT_BRANCH:=$(shell git rev-parse --abbrev-ref HEAD)
endif

#GIT_BRANCH:=$(shell git branch 2> /dev/null | grep '*' | cut -f2 -d' ')
ifndef GIT_HASH
	GIT_HASH:=$(shell git log --format="%h" -n 1 2> /dev/null)
endif

ifndef BUILD_TS
	BUILD_TS:=$(shell date +%FT%T%z)
endif

LDFLAGS = -X 'github.com/Format-C-eft/universal-proxy/internal/config.branch=$(GIT_BRANCH)'\
          -X 'github.com/Format-C-eft/universal-proxy/internal/config.commitHash=$(GIT_HASH)'\
          -X 'github.com/Format-C-eft/universal-proxy/internal/config.timeBuild=$(BUILD_TS)'


.PHONY: .build
.build:
	$(info Building...)
	$(BUILD_ENV_PARAMS) go build -ldflags "$(LDFLAGS)" -o ./bin/universal-proxy ./cmd/universal-proxy

.PHONY: build
build: .bin-deps .build

.PHONY: .bin-deps
.bin-deps:
	$(info Installing binary dependencies...)
