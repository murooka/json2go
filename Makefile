.PHONY: install build-test

install:
	go install github.com/murooka/json2go

build-test: install
	for dir in ./testdata/*; do cd $$dir; pwd; ./build.sh; cd ../../; done
