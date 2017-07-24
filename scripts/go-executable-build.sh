#!/bin/bash

package="github.com/bhklab/PharmacoDB/api/initialize/main.go"

platforms=("linux/amd64" "linux/arm" "linux/386" "windows/amd64" "windows/386" "darwin/amd64" "netbsd/amd64" "openbsd/amd64" "solaris/amd64" "android/arm")

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name='../dist/api/api-'$GOOS'-'$GOARCH
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name $package
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done
