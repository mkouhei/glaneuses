#!/usr/bin/make -f
# -*- makefile -*-

BIN := glaneuses
SRC := *.go
GOPKG := github.com/mkouhei/glaneuses/
GOPATH := $(CURDIR)/_build:$(GOPATH)
export GOPATH


all: precheck clean format test build

precheck:
	@if [ -d .git ]; then \
	set -e; \
	diff -u .git/hooks/pre-commit utils/pre-commit.txt ;\
	[ -x .git/hooks/pre-commit ] ;\
	fi

prebuild:
	go get github.com/kolo/xmlrpc
	go get github.com/bitly/go-simplejson
	go get github.com/miguel-branco/goconfig
	go get -d github.com/PuerkitoBio/goquery
	install -d $(CURDIR)/_build/src/$(GOPKG)
	cp -a $(CURDIR)/*.go $(CURDIR)/_build/src/$(GOPKG)


build: prebuild
	go build -o _build/$(BIN)

build-only:
	go build -o _build/$(BIN)

clean:
	@rm -f _build/$(BIN)

format:
	for src in $(SRC); do \
		gofmt $$src > $$src.tmp ;\
		goimports $$src.tmp > $$src.tmp2 ;\
		mv -f $$src.tmp2 $$src ;\
		rm -f $$src.tmp ;\
	done


test: prebuild
	go test -v -cover
