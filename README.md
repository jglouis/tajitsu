[![Build Status](https://travis-ci.org/jglouis/tajitsu.svg?branch=master)](https://travis-ci.org/jglouis/tajitsu)

# Tajitsu

## Install

`go get github.com/jglouis/tajitsu`

## Test

to run the tests: `go test ./...`

## Update data files

After any update of the data file, must run `go-bindata -o asset/bindata.go -pkg asset data/` within the source folder of the project. 