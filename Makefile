PROJECT  = alfred-chrome-history
TESTARGS ?= -v

dist: build
	(cd build && zip -r "../$(PROJECT).alfredworkflow" .)
.PHONY: dist

build:
	go build -o build/$(PROJECT) -ldflags="-s -w"
	cp _workflow/* build/
.PHONY: build

test:
	go test ./... $(TESTARGS)
.PHONY: test
