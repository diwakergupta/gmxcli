#!/bin/bash

LDFLAGS="-s -w -X main.version=${TRAVIS_TAG:-TRAVIS_COMMIT}"
LDFLAGS+=" -X github.com/diwakergupta/gmxcli/cmd.clientID=${GMXCLI_CLIENT_ID}"
LDFLAGS+=" -X github.com/diwakergupta/gmxcli/cmd.clientSecret=${GMXCLI_CLIENT_SECRET}"

if [[ "$TRAVIS" == "true" ]]; then
	platforms=("darwin/amd64" "linux/amd64" "windows/amd64")
  mkdir -p dist
	for platform in "${platforms[@]}"
	do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=dist/gmxcli'-'$GOOS'-'$GOARCH
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi  

    env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="${LDFLAGS}" -o $output_name $package
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
	done
else
	go build -ldflags="${LDFLAGS}"
fi
