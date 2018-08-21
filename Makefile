.PHONY: install build-test

install:
	go install github.com/murooka/json2go

build-test:
	json2go --package main --typename Person --varname PersonList --output testdata/person.go testdata/persons.json
