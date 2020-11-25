
#Version := $(shell git describe --tags --dirty)
Version := dev
GitCommit := $(shell git rev-parse HEAD)
DIR := $(shell pwd)
LDFLAGS := "-s -w -X main.Version=$(Version) -X main.GitCommit=$(GitCommit)"


.PHONY: build
build:
	CGO_ENABLED=0 go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/skeleton cmd/*.go

.PHONY: dist
dist:
	CGO_ENABLED=0 GOOS=linux go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/skeleton-linux cmd/*.go
	CGO_ENABLED=0 GOOS=darwin go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/skeleton-darwin cmd/*.go	
	CGO_ENABLED=0 GOOS=windows go build -ldflags $(LDFLAGS) -a -installsuffix cgo -o bin/skeleton.exe cmd/*.go	

.PHONY: install
install:
	rm -fr /usr/local/bin/skeleton && ln -s $(DIR)/bin/skeleton /usr/local/bin/skeleton
