# go-diary

go-diary analys and format directory, and rewrite README.md which is written structure directory.


## Usage

```sh
$ diary -h
Usage of diary:
  -copyDir string
    	Format directory. 
    	When this option is difference from 'dir', all file will copy to 'copyDir'.
  -dir string
    	Analysis directory. (default ".")
  -file string
    	Rewrite file.
    	Default value is './README.md.
    	In org mode value is ./README.org
  -org
    	Use org template.
  -tmpl string
    	Parse template file.
    	Default is template/top.template.md.
    	In org mode value is template/org.template.org.
  -v	Output verbose.
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

$ diary
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


### download binary

- [mac 64bit](./build/darwin-amd64/diary)
- [windows 64bit](./build/windows-amd64/diary)
- [linux 64bit](./build/linux-amd64/diary)

## Lisesnce

MIT

## Author

komem3

