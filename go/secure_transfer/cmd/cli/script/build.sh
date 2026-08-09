#!/bin/bash

from_path=$(pwd) # record current path

# change dir, for run this script from anywhere
cd "$(dirname "$0")" || exit 1
cd .. || exit 1

# clean history build folder
if [ -d "./build/" ]; then
  rm -rf ./build/*
fi

mkdir -p "./build"

# Compile multi-platform executables

  compile_exec() {
    local platform="$1"

    IFS='/' read -r goos goarch <<< "$platform"

    local fileName="transfer_cli_${goos}_${goarch}"
    if [ "$goos" = "windows" ]; then
      fileName="$fileName.exe"
    fi
    local filePath="./build/$fileName"

    GOOS="$goos" GOARCH="$goarch" go build -o "$filePath"

    sha1sum "$filePath" | cut -d" " -f1 > "$filePath.sha1"
  }

go mod tidy # 这么做其实会下载全部依赖，所以该脚本仅供参考

compile_exec "windows/amd64"
compile_exec "linux/amd64"
compile_exec "linux/arm64"

cp "./script/manual.md" "./build/manual_cli.md"
sed -i "3i\> Build Time: \\$(date)\n> Go Version: \\$(go version)\n" "./build/manual_cli.md"

# back to from path
cd "$from_path" || exit 1
