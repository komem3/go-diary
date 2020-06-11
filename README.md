# Go-diary

Go-diary is a tool for managing diaries.

## Contents

- [Go-diary](#go-diary)
  - [Usage](#usage)
    - [Initialize](#initialize)
    - [New](#new)
    - [Format](#format)
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
  new         Generate new diary

Flags:
  -h, --help      help for diary

  Use "diary [command] --help" for more information about a command.
```

### Initialize

```sh
$ diary init -h
Init command make template directory.
You need to run this command before running other command.

Usage:
diary init [flags]

Flags:
-d, --dir string   Created template directory path (default ".")
-h, --help         help for init
--v            Output verbose.
```

### New

```sh
$ diary new -h
New command create new today diary from template file.

Usage:
diary new [flags]

Flags:
-d, --date string     Date of making diary.
Format: YYYY/MM/dd(2010/01/31) or today(t) or yesterday(y) or tomorrow(tm).
(default "today")
--dir string      Destination directory. (default ".")
-f, --format string   File name format.
Refer to https://golang.org/src/time/format.go (default "20060102.md")
-h, --help            help for new
--tmpl string     Parse template file. (default "template/diary.template.md")
```

### Format

```sh
$ diary format -h
Format command analys and format directory.
After format directory, it write directory structure to target file.

Usage:
diary format [flags]

Flags:
--copyDir string   Format directory.
When this option is difference from 'dir', all file will copy to 'copyDir'.
-d, --dir string       Analysis directory. (default ".")
-f, --file string      Write file. (default "./README.md")
-h, --help             help for format
--tmpl string      Parse template file. (default "template/top.template.md")
--v                Output verbose.
```



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
