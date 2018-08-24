#!/bin/bash

json2go --package main --typename Example --varname Examples --structure map --output example.go example.json
