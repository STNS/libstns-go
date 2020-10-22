VERSION = $(shell cat version)
REVISION := $(shell git rev-parse --short HEAD)
INFO_COLOR=\033[1;34m
RESET=\033[0m
BOLD=\033[1m
TEST ?= $(shell go list ./... | grep -v -e vendor -e keys -e tmp)

GOVERSION=$(shell go version)
GO ?= GO111MODULE=on go

TESTCONFIG="misc/test.conf"

DIST ?= unknown
PREFIX=/usr
BINDIR=$(PREFIX)/sbin
SOURCES=Makefile go.mod go.sum version cmd cache_stnsd main.go package/
BUILD=tmp/bin
UNAME_S := $(shell uname -s)
.DEFAULT_GOAL := build

.PHONY: build
## build: build the nke
build:
	$(GO) build -o $(BUILD)/cache-stnsd -ldflags "-X github.com/STNS/cache-stnsd/cmd.version=$(VERSION)"

.PHONY: release
## release: release nke (tagging and exec goreleaser)
release:
	goreleaser --rm-dist

.PHONY: bump
bump:
	git semv minor --bump
	git tag | tail -n1 | sed 's/v//g' > version

.PHONY: releasedeps
releasedeps: git-semv goreleaser

.PHONY: git-semv
git-semv:
	brew tap linyows/git-semv
	brew install git-semv


.PHONY: goreleaser
goreleaser:
	brew install goreleaser/tap/goreleaser
	brew install goreleaser

.PHONY: test
test:
	@echo "$(INFO_COLOR)==> $(RESET)$(BOLD)Testing$(RESET)"
	$(GO) test -v $(TEST) -timeout=30s -parallel=4
	$(GO) test -race $(TEST)

.PHONY: github_release
github_release: ## Create some distribution packages
	ghr -u STNS --replace v$(VERSION) builds/
