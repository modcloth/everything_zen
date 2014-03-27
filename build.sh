#!/bin/sh -e

if [ ! -h src/github.com/modcloth/everything_zen ]; then
  mkdir -p src/github.com/modcloth
  ln -s ../../.. src/github.com/modcloth/everything_zen
fi

export GOBIN=${PWD}/bin
export GOPATH=${PWD}:$(godep path)

go install github.com/modcloth/everything_zen
