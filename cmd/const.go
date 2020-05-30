package cmd

import "path/filepath"

const tmplDir = "template"

func mdTmp() string {
	return filepath.Join(tmplDir, "top.template.md")
}

func orgTmp() string {
	return filepath.Join(tmplDir, "org.template.org")
}
