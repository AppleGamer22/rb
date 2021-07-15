# done -> CompletionStatus
# json -> target
# target usb, source smb

VERSION:=$(shell cat ./VERSION)
SOURCES:=./pkg/rb
# EXSOURCES:=./examples/maxp
BIN:=./bin
BASE:=rb
LINUX=$(BASE)_linux_amd64
DARWIN=$(BASE)_darwin_amd64
WINDOWS=$(BASE)_win_amd64.exe



all: build test build-example test-integ benchmark

build:
		@echo "Building $(VERSION)"
		go build -v $(SOURCES)

test:
		@echo
		@echo "Testing $(VERSION)"
		go test -v -cover $(SOURCES)

benchmark:
		@echo
		@echo "Benchmarking $(VERSION)"
		go test -run=XXX -bench=. $(SOURCES)

linux: $(LINUX)

darwin: $(DARWIN)

windows: $(WINDOWS)

$(LINUX):
	env GOOS=linux GOARCH=amd64 go build -v -o $(BIN)/$(LINUX) $(EXSOURCES)

$(DARWIN):
	env GOOS=darwin GOARCH=amd64 go build -v -o $(BIN)/$(DARWIN) $(EXSOURCES)

$(WINDOWS):
	env GOOS=windows GOARCH=amd64 go build  -v -o $(BIN)/$(WINDOWS) $(EXSOURCES)
	
build-example: linux darwin windows
	@echo
	@echo "Example application build $(VERSION) complete"

test-integ:
		@echo
		@echo "Integration test $(VERSION)"
		go test $(EXSOURCES)

.PHONY: build test benchmark build-example test-integ