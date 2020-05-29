# go-diary

go-diary analys and format directory, and rewrite README.md which is written structure directory.


## Usage

```sh
$ diary -h
Diary is a CLI libray for managing your diary.
This application can format your diary directory, and make index file.

Usage:
  diary [command]

Available Commands:
  format      Format directory
  help        Help about any command
  init        Initialize directory

Flags:
  -h, --help      help for diary
  -v, --version   version for diary

Use "diary [command] --help" for more information about a command.
```


## Target file

Target diary file is `YYYYMMDD*.*`.
- good example: *20200102.md*, *19990912_test.md*, *19990912.txt*
- bad example: *diary_20111102.md*, *202011.md*


## Sample case

```sh
$ pwd
${GOPATH}/src/github.com/komem3/diary/sample

$ tree
.
├── 20021201.md
├── 20200101.md
├── 20200102.txt
└── 20200202.org

$ diary init
$ ls
20021201.md  20200101.md  20200102.txt  20200202.org  template

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
    ├── org.template.org
    └── top.template.md

$ cat README.md
# diary record

## 2020

<details>
<summary>01</summary>
<ul>
<li><a href="./2020/01/20200101.md">20200101.md</a></li>
<li><a href="./2020/01/20200102.txt">20200102.txt</a></li><ul></details>

<details>
<summary>02</summary>
<ul>
<li><a href="./2020/02/20200202.org">20200202.org</a></li><ul></details>

## 2002

<details>
<summary>12</summary>
<ul>
<li><a href="./2002/12/20021201.md">20021201.md</a></li><ul></details>

```


## Install

### Install with go tool

```sh
go get github.com/komem3/diary/cmd/diary
```

## Lisesnce

MIT

## Author

komem3

