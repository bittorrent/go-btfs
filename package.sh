#!/bin/bash

echo "Automatic update script begin!"

version=`cat version.go | tail -n +7 | head -n +1 | awk '{print $4}'`

echo "==========>Get btfs version [${version}]"

os=( darwin linux windows )
arch=( amd64 386 arm64 arm )

# Get the version for fs-repo-migrations
curl -fsSL -O https://github.com/TRON-US/btfs-distributions/raw/master/fs-repo-migrations/versions
if [ $? -eq 0 ]
then
    VERSION=$(cat versions | sort -V | tail -n 1)
    rm -f versions
else
    echo "Download of fs-repo-migrations version file failed, confirm your internet connection to GitHub is working and rerun this script."
    exit 1
fi

if ! command -v unzip > /dev/null; then apt-get update; apt-get -y install unzip; fi

for goos in ${os[@]}
do
    for goarch in ${arch[@]}
    do
        if [[ ( ${goos} = "windows" && ( ${goarch} = "arm64" || ${goarch} = "arm" ) ) || ( ${goos} = "darwin" && ( ${goarch} = "arm"  || ${goarch} = "386" ) )]]; then continue; fi
        echo "=============>OS: [${goos}] ARCH: [${goarch}] automatic compiler begin."
        ext=""
        if [[ ${goos} = "windows" ]]; then
            ext=".exe"
        fi
        rm -f ../btfs-binary-releases/${goos}/${goarch}/*
        GOOS=${goos} GOARCH=${goarch} GOFLAGS=-buildvcs=false make build
        GOOS=${goos} GOARCH=${goarch} GOFLAGS=-buildvcs=false go build -o update-${goos}-${goarch}${ext} autoupdate/main.go
        md5=`openssl md5 ./cmd/btfs/btfs | awk '{print $2}'`
        mv ./cmd/btfs/btfs ../btfs-binary-releases/${goos}/${goarch}/btfs-${goos}-${goarch}${ext}
        mv update-${goos}-${goarch}${ext} ../btfs-binary-releases/${goos}/${goarch}/update-${goos}-${goarch}${ext}
        if [[ ${goos} != "windows" ]]
        then
            if [[ ${goos} = "darwin" ]] && [[ ${goarch} = "arm64" ]]; then
                echo "Warning: Skipping fs-repo-migrations for darwin/arm64"
            else
                URL="https://github.com/TRON-US/btfs-distributions/raw/master/fs-repo-migrations/${VERSION}/fs-repo-migrations_${VERSION}_${goos}-${goarch}.tar.gz"
                wget -q "$URL"
                tar -xzf fs-repo-migrations_${VERSION}_${goos}-${goarch}.tar.gz
                mv fs-repo-migrations/fs-repo-migrations ../btfs-binary-releases/${goos}/${goarch}/fs-repo-migrations-${goos}-${goarch}
                rm -f fs-repo-migrations_${VERSION}_${goos}-${goarch}.tar.gz
                rm -rf fs-repo-migrations
            fi
        else
            URL="https://github.com/TRON-US/btfs-distributions/raw/master/fs-repo-migrations/${VERSION}/fs-repo-migrations_${VERSION}_${goos}-${goarch}.zip"
            wget -q "$URL"
            unzip -q fs-repo-migrations_${VERSION}_${goos}-${goarch}.zip
            mv fs-repo-migrations/fs-repo-migrations.exe ../btfs-binary-releases/${goos}/${goarch}/fs-repo-migrations-${goos}-${goarch}.exe
            rm -f fs-repo-migrations_${VERSION}_${goos}-${goarch}.zip
            rm -rf fs-repo-migrations
        fi
        echo -e "version: ${version: 1: ${#version}-2}\nmd5: ${md5}\nautoupdateFlg: true\nsleepTimeSeconds: 86400\nbeginNumber: 0\nendNumber: 100" > ../btfs-binary-releases/${goos}/${goarch}/config_${goos}_${goarch}.yaml
        cd ../btfs-binary-releases/${goos}/${goarch}
        tar -cvf btfs-${goos}-${goarch}.tar config_${goos}_${goarch}.yaml btfs-${goos}-${goarch}${ext}
        cd -
        echo "=============>OS: [${goos}] ARCH: [${goarch}] automatic compiler success."
    done
done

echo "Automatic update script success!"