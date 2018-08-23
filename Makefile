.PHONY: install build-test

install:
	go install github.com/murooka/json2go

build-test: install
	for dir in ./testdata/*; do cd $$dir; json2go --package main --typename Example --varname Examples --output example.go example.json; cd ../../; done
