VETARGS?=-all
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
BIN_NAME?=kube-etcd-helper
GO_FILES?=$(shell find . -name '*.go')
GOVERSION=$(shell go version)
THIS_GOOS=$(word 1,$(subst /, ,$(lastword $(GOVERSION))))
THIS_GOARCH=$(word 2,$(subst /, ,$(lastword $(GOVERSION))))
GOOS=$(THIS_GOOS)
GOARCH=$(THIS_GOARCH)
BUILD_LDFLAGS = "-s -w"
BUILD_TARGETS= \
	build-linux-arm64 \
	build-linux-arm \
	build-linux-amd64 \
	build-linux-386 \
	build-darwin-amd64 \
	build-darwin-386 \
	build-windows-amd64 \
	build-windows-386

.PHONY: default
default: test build

.PHONY: run
run:
	go run $(CURDIR)/*.go $(ARGS)

.PHONY: clean
clean:
	rm -Rf bin/*

.PHONY: tools
tools:
	go get -u github.com/golang/dep/cmd/dep
	go get -v github.com/alecthomas/gometalinter
	gometalinter --install

.PHONY: build $(BUILD_TARGETS) bin/$(GOOS)/$(GOARCH)/peco$(SUFFIX)


.PHONY: build
build: bin/kube-etcd-helper_$(GOOS)_$(GOARCH)$(SUFFIX)

build-all: $(BUILD_TARGETS)

build-windows-amd64:
	@$(MAKE) build GOOS=windows GOARCH=amd64 SUFFIX=.exe

build-windows-386:
	@$(MAKE) build GOOS=windows GOARCH=386 SUFFIX=.exe

build-linux-amd64:
	@$(MAKE) build GOOS=linux GOARCH=amd64

build-linux-arm:
	@$(MAKE) build GOOS=linux GOARCH=arm

build-linux-arm64:
	@$(MAKE) build GOOS=linux GOARCH=arm64

build-linux-386:
	@$(MAKE) build GOOS=linux GOARCH=386

build-darwin-amd64:
	@$(MAKE) build GOOS=darwin GOARCH=amd64

build-darwin-386:
	@$(MAKE) build GOOS=darwin GOARCH=386

bin/kube-etcd-helper_$(GOOS)_$(GOARCH)$(SUFFIX): $(GO_FILES)
	CGO_ENABLED=0 go build -ldflags $(BUILD_LDFLAGS) -o bin/kube-etcd-helper_$(GOOS)_$(GOARCH)$(SUFFIX) *.go

.PHONY: test
test: lint
	go test ./... $(TESTARGS) -v -timeout=30m -parallel=4 ;

.PHONY: lint
lint: fmt
	gometalinter --vendor --skip=vendor/ --cyclo-over=15 --disable=gas --disable=maligned --deadline=2m ./...
	@echo

.PHONY: fmt
fmt:
	gofmt -s -l -w $(GOFMT_FILES)

