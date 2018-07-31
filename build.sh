#!/bin/bash

LDFLAGS="-s -w -X main.version=${TRAVIS_TAG:-TRAVIS_COMMIT}"
LDFLAGS+=" -X github.com/diwakergupta/gmxcli/cmd.clientID=${GMXCLI_CLIENT_ID}"
LDFLAGS+=" -X github.com/diwakergupta/gmxcli/cmd.clientSecret=${GMXCLI_CLIENT_SECRET}"
go build -ldflags="${LDFLAGS}"
