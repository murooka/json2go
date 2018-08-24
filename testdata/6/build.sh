#!/bin/bash

json2go --package main --typename Example --varname Examples --root /data --structure map,map --output example.go example1.json example2.json
