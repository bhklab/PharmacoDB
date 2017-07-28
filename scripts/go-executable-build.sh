#!/bin/bash

echo "Creating latest builds for all defined platforms ..."

package="../api/initialize/api.go"

platforms=("linux/amd64"
           "linux/arm"
           "linux/386"
           "windows/amd64"
           "windows/386"
           "darwin/amd64"
           "netbsd/amd64"
           "openbsd/amd64")

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

echo "Finished successfully!"
