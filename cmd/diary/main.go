package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/komem3/diary"
	_ "github.com/komem3/diary/statik"
	"github.com/rakyll/statik/fs"
)

const tmpDir = "template"

var (
	mdTmp  = filepath.Join(tmpDir, "top.template.md")
	orgTmp = filepath.Join(tmpDir, "org.template.org")
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
	logger := diary.NewLogger(verbose)

	logger.Debug(
		"msg", "initialize template",
	)
	if err := initialize(logger); err != nil {
		logger.Error(
			"when", "initialize",
			"err", err,
		)
		return
	}

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
	}
}

func initialize(l diary.Logger) error {
	err := os.Mkdir(tmpDir, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}

	statikFS, err := fs.New()
	if err != nil {
		return err
	}

	createTmp := func(tmpName string) error {
		r, err := statikFS.Open("/" + filepath.Base(tmpName))
		if err != nil {
			return err
		}
		defer diary.CloseWithErrLog(l, r)
		contents, err := ioutil.ReadAll(r)
		if err != nil {
			return err
		}

		file, err := os.OpenFile(tmpName, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer diary.CloseWithErrLog(l, file)

		_, err = file.Write(contents)
		if err != nil {
			return err
		}
		return nil
	}

	l.Debug(
		"msg", "create "+mdTmp,
	)
	if err = createTmp(mdTmp); err != nil {
		return err
	}
	l.Debug(
		"msg", "create "+orgTmp,
	)
	if err = createTmp(orgTmp); err != nil {
		return err
	}
	return nil
}
