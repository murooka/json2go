#!/bin/bash

json2go --package main --typename UserAgent --varname UserAgents --root /userAgents --structure map,map --output example.go example1.json example2.json
