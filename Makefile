# done -> CompletionStatus
# json -> target
# target usb, source smb

VERSION:=1.0.0
BIN:=./bin
BASE:=rb_$(VERSION)
LINUX=$(BASE)_linux_amd64
DARWIN=$(BASE)_darwin_amd64
WINDOWS=$(BASE)_windows_amd64.exe

all: test build

build: $(LINUX) $(DARWIN) $(WINDOWS)

test:
	@echo "Testing $(VERSION) internal"
	go test -v -cover ./internal/*
	@echo "Testing $(VERSION) cmd"
	go test -v -cover ./cmd/*

linux: $(LINUX)

darwin: $(DARWIN)

windows: $(WINDOWS)

$(LINUX):
	env GOOS=linux GOARCH=amd64 go build -v -o $(BIN)/$(LINUX)

$(DARWIN):
	env GOOS=darwin GOARCH=amd64 go build -v -o $(BIN)/$(DARWIN)

$(WINDOWS):
	env GOOS=windows GOARCH=amd64 go build  -v -o $(BIN)/$(WINDOWS)

.PHONY: build test