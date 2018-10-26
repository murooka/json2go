.PHONY: help install build-test test

## show this message
help:
	@make2help

## install json2go
install:
	go install github.com/murooka/json2go

## build test cases
build-test: install
	for dir in ./testdata/*; do cd $$dir; pwd; ./build.sh; cd ../../; done

## run test
test:
	go test -v ./...
