#!/usr/bin/make -f
# -*- makefile -*-

BIN := glaneuses
SRC := *.go
GOPKG := github.com/mkouhei/glaneuses/
GOPATH := $(CURDIR)/_build:$(GOPATH)
export GOPATH


#all: clean format test build
all: clean format build

prebuild:
	go get github.com/kolo/xmlrpc
	go get github.com/bitly/go-simplejson
	go get github.com/msbranco/goconfig
	go get -d github.com/PuerkitoBio/goquery
	install -d $(CURDIR)/_build/src/$(GOPKG)
	cp -a $(CURDIR)/*.go $(CURDIR)/_build/src/$(GOPKG)


build: prebuild
	go build -o _build/$(BIN)

build-only:
	go build -o _build/$(BIN)

clean:
	#@rm -rf _build/glaneuses
	@rm -rf _build/


format:
	for src in $(SRC); do \
		gofmt $$src > $$src.tmp ;\
		mv -f $$src.tmp $$src ;\
	done


test: prebuild
	go test
