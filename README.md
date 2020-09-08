# Go-diary

[![CircleCI](https://circleci.com/gh/komem3/go-diary.svg?style=svg&circle-token=e9c905abe7b4522efa323463121ffa00ad8acc0d)](https://app.circleci.com/pipelines/github/komem3/go-diary)


Go-diary is a tool for managing diaries.

## Contents

- [Go-diary](#go-diary)
  - [Usage](#usage)
  - [Install](#install)
    - [Install with go tool](#install-with-go-tool)
  - [Target File](#target-file)
  - [Sample Case](#sample-case)
  - [Template Variables](#template-variables)
    - [`top.template.md`](#toptemplatemd)
    - [`diary.template.md`](#diarytemplatemd)
  - [Licence](#licence)
  - [Author](#author)

## Usage

- [diary](./cmd/doc/diary.md)
  - [init](./cmd/doc/diary_init.md)
  - [new](./cmd/doc/diary_new.md)
  - [format](./cmd/doc/diary_format.md)

## Install

### Install with go tool

```sh
go get github.com/komem3/go-diary/cmd/diary
```

## Target File

Target diary file is `YYYYMMDD*.*`.

- good example: `20200102.md`, `19990912_test.md`, `19990912.txt`
- bad example: `diary_20111102.md`, `202011.md`

## Sample Case

```sh
$ pwd
${GOPATH}/src/github.com/komem3/diary/sample

$ tree
.
├── 20021201.md
├── 20200101.md
└── 20200102.txt

$ diary init
$ ls
20021201.md  20200101.md  20200102.txt template

$ diary format
$ tree
.
├── 2002
│   └── 12
│       └── 20021201.md
├── 2020
│   ├── 01
│   │   ├── 20200101.md
│   │   └── 20200102.txt
│   └── 02
│       └── 20200202.org
├── README.md
└── template
    ├── diary.template.md
    └── top.template.md

$ diary new
$ ls
2002  2020  20200531.md  README.md  template
$ cat 20200531.md
2020/05/31 (Sunday)
```

## Template Variables

### `top.template.md`

- Base: Base path.
- Years: Slice of year element.
  - Year
  - Months: Slice of month element.
    - Month
    - Days: Slice of day element.
      - Day
      - Path: Path from base.

### `diary.template.md`

- Year
- Month
- Day
- Weekday

## Licence

MIT

## Author

komem3
