MAKEFILE_PATH := $(GOPATH)\src\github.com\ixoja\shorten
BIN_PATH := $(MAKEFILE_PATH)\bin
PATH := $(MAKEFILE_PATH):$(PATH)


# Basic go commands
GOCMD     = go
GOBUILD   = $(GOCMD) build
GORUN     = $(GOCMD) run
GOCLEAN   = $(GOCMD) clean
GOTEST    = $(GOCMD) test -race -v -count=1
GOLINT    = golangci-lint

# Binary output name
BINARY = shorten.exe

#
VENDOR_DIR           = ./vendor

#
PKGS = $(shell go list ./... | grep -v /vendor | grep -v grpc)

# Colors
GREEN_COLOR   = "\033[0;32m"
PURPLE_COLOR  = "\033[0;35m"
DEFAULT_COLOR = "\033[m"

# Tests
TEST_STRING = "TEST"

.PHONY: all help clean dep test build run docker swagger-clean swagger

default: clean swagger build


help:
	@echo 'Usage: make <TARGETS> ... <OPTIONS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@echo '    build              Compile packages and dependencies.'
	@echo '    clean              Remove binary.'
	@echo '    dep                Download and install build time dependencies.'
	@echo '    docker             Run docker buld and run.'
	@echo '    help               Show this help screen.'
	@echo '    lint               Run golangci-lint on package sources.'
	@echo '    run                Compile and run Go program.'
	@echo '    test               Run unit tests.'
	@echo '    swagger            Generate swagger models and server.'
	@echo ''
	@echo 'Targets run by default are: clean swagger build'
	@echo ''

clean:
	@echo -e $(GREEN_COLOR)[clean]$(DEFAULT_COLOR)
	@$(GOCLEAN)
	@if [ -f $(BIN_PATH)/$(BINARY) ] ; then rm $(BIN_PATH)/$(BINARY) ; fi

dep:
	@echo -e $(GREEN_COLOR)[DEP CLEAN AND ENSURE]$(DEFAULT_COLOR)
ifneq ("$(wildcard ./Gopkg.lock)","")
	@echo "rm -rf ./Gopkg.lock"
	@rm -rf ./Gopkg.lock
endif
ifneq ("$(wildcard ./vendor)","")
	@echo "rm -rf ./vendor"
	@rm -rf ./vendor
endif
	@$(DEPCMD) ensure -v

lint:
	@echo -e $(GREEN_COLOR)[$(TEST_STRING)]$(DEFAULT_COLOR)
	$(GOLINT) run --enable=golint --enable=stylecheck --enable=gosec --enable=interfacer --enable=unconvert \
	--enable=dupl --enable=goconst --enable=gocyclo --enable=gofmt --enable=maligned --enable=depguard \
	--enable=misspell --enable=lll --enable=unparam --enable=nakedret --enable=prealloc --enable=scopelint \
	--enable=gocritic --enable=gochecknoinits --enable=gochecknoglobals \
	--skip-dirs=restapi \

test:
	@echo -e $(GREEN_COLOR)[$(TEST_STRING)]$(DEFAULT_COLOR)
	@$(GOTEST) $(PKGS)

build:
	@echo -e $(GREEN_COLOR)[build]$(DEFAULT_COLOR)
	@$(GOBUILD) -v -o $(BIN_PATH)/$(BINARY)

run webserver: build
	@echo -e $(GREEN_COLOR)[run web_server]$(DEFAULT_COLOR)
	$(BIN_PATH)\$(BINARY) -mode "webserver" -html $(MAKEFILE_PATH)\internal\html -port 8080 -api_url "localhost:9000"

run apiserver: build
	@echo -e $(GREEN_COLOR)[run api_server]$(DEFAULT_COLOR)
	$(BIN_PATH)\$(BINARY) -mode "apiserver" -api_address ":9000"

docker: build
	@echo -e $(GREEN_COLOR)[DOCKER]$(DEFAULT_COLOR)
	docker-compose up
