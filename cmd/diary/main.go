package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/komem3/diary"
	_ "github.com/komem3/diary/statik"
)

const tmplDir = "template"

var (
	mdTmp  = filepath.Join(tmplDir, "top.template.md")
	orgTmp = filepath.Join(tmplDir, "org.template.org")
)

func main() {
	var (
		from     string
		to       string
		file     string
		tempFile string
		orgMode  bool
		verbose  bool
	)
	{
		flag.StringVar(&from, "dir", ".", "Analysis directory.")
		flag.StringVar(&to, "copyDir", "", "Format directory. \nWhen this option is difference from 'dir', all file will copy to 'copyDir'.")
		flag.StringVar(&file, "file", "", "Rewrite file.\nDefault value is './README.md.\nIn org mode value is ./README.org")
		flag.StringVar(&tempFile, "tmpl", "",
			"Parse template file.\n"+
				"Default is "+mdTmp+".\n"+
				"In org mode value is "+orgTmp+".")
		flag.BoolVar(
			&orgMode,
			"org",
			false,
			"Use org template.")
		flag.BoolVar(&verbose, "v", false, "Output verbose.")

		flag.Parse()
	}

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	logger := diary.NewLogger(verbose)
	if orgMode {
		if file == "" {
			file = "./README.org"
		}
		if tempFile == "" {
			tempFile = orgTmp
		}
	} else {
		if file == "" {
			file = "./README.md"
		}
		if tempFile == "" {
			tempFile = mdTmp
		}
	}

	switch flag.Arg(0) {
	case "init":
		logger.Debug(
			"msg", "initialize template",
		)
		if err := diary.Initialize(logger, ".", tmplDir, []string{mdTmp, orgTmp}); err != nil {
			logger.Error(
				"when", "initialize",
				"err", err,
			)
			os.Exit(1)
		}
	case "format":
		dGen, err := diary.NewDiaryGenerator(file, tempFile, from, to, orgMode, logger)
		if err != nil {
			logger.Error(
				"when", "generate generator",
				"err", err.Error(),
			)
			return
		}
		logger.Debug(
			"msg", "analys start",
			"dir", from,
		)
		fMap := diary.ParseFileMap(from)
		elem := diary.Map2Elem(fMap)
		dGen.FormatDir(fMap).WriteDirTree(elem)

		if dGen.Err != nil {
			logger.Error(
				"when", "format and write dir tree",
				"err", dGen.Err.Error(),
			)
			os.Exit(1)
		}
	default:
		flag.Usage()
		fmt.Printf("usage?%v\n", flag.Arg(1))
		os.Exit(1)
	}
	os.Exit(0)
}
