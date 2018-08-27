#!/bin/bash

json2go --package main --typename Example --varname Examples --structure map,slice --output example.go example1.json example2.json
