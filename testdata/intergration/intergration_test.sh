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

function removeWord() {
    cat - | sed -e 's/generated //g'
}

outputCheck $(diary new | removeWord)
outputCheck $(diary new -d 1999/1/22 | removeWord) 19990122.md
outputCheck $(diary new -d y | removeWord)
outputCheck $(diary new -d tm | removeWord)

echo "test" > template/copy.template.md
outputCheck $(DIARY_NEW_TEMPLATE=template/copy.template.md diary new -d 2001/1/1 | removeWord) 20010101.md
outputCheck $(cat 20010101.md) "test"

echo "flag" > template/flag.template.md
outputCheck $(diary new -d 2001/2/1 --tmpl template/flag.template.md | removeWord) 20010201.md
outputCheck $(cat 20010201.md) "flag"

echo "* ext" > template/ext.template.org
outputCheck $(diary new -d 2001/2/2 --format "2006_01_02" --tmpl template/ext.template.org | removeWord) 2001_02_02.org
outputCheck $(cat 2001_02_02.org) "* ext"

function removeWord() {
    cat - | sed -e 's/write index to //g'
}

outputCheck $(diary format | removeWord) README.md
outputCheck $(diary format --copyDir copy --file copy.md --yearSort asc --monthSort desc --daySort desc | removeWord) copy/copy.md

outputCheck $(DIARY_INDEX_FILE=index.txt DIARY_INDEX_TEMPLATE=template/copy.template.md \
                              diary format | removeWord) index.txt
outputCheck $(cat index.txt) "test"

outputCheck $(diary format --file index.out --tmpl template/flag.template.md | removeWord) index.out
outputCheck $(cat index.out) "flag"

outputCheck $(diary format --file index --tmpl template/ext.template.org | removeWord) index.org
outputCheck $(cat index.org) "* ext"

dir=${PWD##*/}
if [ $dir != "mydiary" ]; then
    echo "${dir} directory is not mydiary directory." 1>&2
    exit 1
fi

echo 'Intergration test is success!'
