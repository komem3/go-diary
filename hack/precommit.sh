#!/bin/bash

# Prettify all selected files
FILES=$(git diff --cached --name-only --diff-filter=ACMR "*.md" "*.yml" "*.json" | sed 's| |\\ |g')
if [ -n "$FILES" ];then
    echo "$FILES" | xargs ./node_modules/.bin/prettier --write
    echo "$FILES" | xargs git add
fi

# Go formatter
FILES=$(git diff --cached --name-only --diff-filter=ACMR "*.go" | sed 's| |\\ |g')
if [ -n "$FILES" ];then
    echo "$FILES" | xargs go fmt
    echo "$FILES" | xargs goimports
    echo "$FILES" | xargs git add
fi

FILES=$(git diff --cached --name-only --diff-filter=ACMR "go.mod" "go.sum" | sed 's| |\\ |g')
if [ -n "$FILES" ];then
    go mod tidy
    git add go.mod go.sum
fi

exit 0
