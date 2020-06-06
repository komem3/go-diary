#!/bin/bash

# set -o noexec
set -o errexit

dir=${PWD##*/}
if [ $dir != "intergration" ]; then
    echo "${dir} directory is not intergration directory." 1>&2
    exit 1
fi

go get -u github.com/komem3/go-diary/cmd/diary

mkdir mydiary
cd mydiary

trap 'cd ../ && rm -r mydiary' EXIT 1 2 3 15

diary -h
diary -v

diary init

diary new
diary new -d 1999/1/22
diary new -d y
diary new -d tm

diary format
diary format --copyDir copy --file copy.md

dir=${PWD##*/}
if [ $dir != "mydiary" ]; then
    echo "${dir} directory is not mydiary directory." 1>&2
    exit 1
fi

tree

echo 'Intergration test is success!'
