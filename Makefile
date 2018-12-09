.PHONY: usage build run test

OK_COLOR=\033[32;01m
NO_COLOR=\033[0m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

GO = go
PKGS = $(shell $(GO) list ./...)
BIN = gitlabels-cli

ECHOFLAGS ?=

BUILDOS ?= linux
BUILDARCH ?= amd64
BUILDENVS ?= CGO_ENABLED=0 GOOS=$(BUILDOS) GOARCH=$(BUILDARCH)
BUILDFLAGS ?= -a -installsuffix cgo --ldflags '-extldflags "-lm -lstdc++ -static"'

usage: Makefile
	@echo $(ECHOFLAGS) "to use make call:"
	@echo $(ECHOFLAGS) "    make <action>"
	@echo $(ECHOFLAGS) ""
	@echo $(ECHOFLAGS) "list of available actions:"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'

## build: build gitlabels-cli
build: test
	@echo $(ECHOFLAGS) "$(OK_COLOR)==> Building binary ($(BUILDOS)/$(BUILDARCH)/$(BIN))...$(NO_COLOR)"
	@$(BUILDENVS) $(GO) build -v $(BUILDFLAGS) -o bin/$(BUILDOS)_$(BUILDARCH)/$(BIN) ./cmd/git-labels-cli

## test: run unit tests
test:
	@echo $(ECHOFLAGS) "$(OK_COLOR)==> Running tests...$(NO_COLOR)"
	@$(GO) test $(PKGS)

## run: run gitlabels-cli
run:
	@echo $(ECHOFLAGS) "$(OK_COLOR)==> Running bin/$(BUILDOS)_$(BUILDARCH)/$(BIN)...$(NO_COLOR)"
	@./bin/$(BUILDOS)_$(BUILDARCH)/$(BIN) $(args)
