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

go mod tidy

# 该脚本仅供参考，fyne有交叉编译问题，只按照开始要求下载的内容，无法在linux上编译出windows可执行程序
# 虽然go编译的二进制程序不需要额外的运行时环境，但是fyne依赖cgo，对系统glibc版本有要求，所以不提供gui可执行程序，仅在文档中附图
fyne package -os linux -icon ./script/icon_256.png --name transfer_gui_linux -release
sha1sum "transfer_gui_linux.tar.xz" | cut -d" " -f1 > "transfer_gui_linux.tar.xz.sha1"
mv "./transfer_gui_linux.tar.xz" "./build/transfer_gui_linux.tar.xz"
mv "./transfer_gui_linux.tar.xz.sha1" "./build/transfer_gui_linux.tar.xz.sha1"

fyne package -os windows -icon ./script/icon_256.png --name transfer_gui_windows.exe --app-id secure.transfer -release
sha1sum "transfer_gui_windows.exe" | cut -d" " -f1 > "transfer_gui_windows.exe.sha1"
mv "./transfer_gui_windows.exe" "./build/transfer_gui_windows.exe"
mv "./transfer_gui_windows.exe.sha1" "./build/transfer_gui_windows.exe.sha1"

# back to from path
cd "$from_path" || exit 1
