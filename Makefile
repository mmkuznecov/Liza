PREFIX=/usr/local/las2

ifndef BUILDHASH
BUILDHASH ?= $(shell git log -1 --pretty='format:%h')
endif

PATH:=$(PWD)/bin:$(PATH)
GOPATH=$(PWD):$(PWD)/vendor
export GOPATH

TEST_FILES = $(shell find src -name '*_test.go')
TEST_PACKAGES = $(shell for f in $(TEST_FILES) ; do dirname $$f | sed -e 's/src\///' ; done | sort -u)

BUILDNUMBER ?= $(shell date +%Y%m%d.%H%M-%S)
VERSION     ?= $(BUILDHASH)
LDFLAGS   = "-X main.version=$(VERSION) -X main.buildNumber=$(BUILDNUMBER)"

RPM_VERSION ?= $(shell date +%Y%m%d.%H%M)
RPM_RELEASE ?= $(shell date +%S)

all: build_las2

build_las2: bin/gb
	bin/gb build -ldflags $(LDFLAGS) las2

clean:
	rm -rf pkg bin

bin/gb:
	go build -o bin/gb github.com/constabulary/gb/cmd/gb
	go build -o bin/gb-vendor github.com/constabulary/gb/cmd/gb-vendor

test:
	go test $(TEST_PACKAGES)
