#!/bin/bash

function outputCheck() {
    if [ -z "$1" ]; then
        echo "no output" 1>&2
        exit 1
    fi
    if [ $# -eq 2 ] && [ "$1" != "$2" ]; then
        echo "$1 should be $2" 1>&2
        exit 1
    fi
}

set -o nounset
set -o errexit

dir=${PWD##*/}
if [ $dir != "intergration" ]; then
    echo "${dir} directory is not intergration directory." 1>&2
    exit 1
fi

mkdir mydiary
cd mydiary
trap 'cd ../ && rm -r mydiary' EXIT 1 2 3 15

outputCheck $(diary -h)

diary init
outputCheck "$(ls template)"

outputCheck $(diary new | sed -e 's/generated //g')
outputCheck $(diary new -d 1999/1/22 | sed -e 's/generated //g') 19990122.md
outputCheck $(diary new -d y | sed -e 's/generated //g')
outputCheck $(diary new -d tm | sed -e 's/generated //g')

diary format
diary format --copyDir copy --file copy.md

dir=${PWD##*/}
if [ $dir != "mydiary" ]; then
    echo "${dir} directory is not mydiary directory." 1>&2
    exit 1
fi

echo 'Intergration test is success!'
