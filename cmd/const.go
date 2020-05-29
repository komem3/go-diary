package cmd

import "path/filepath"

const tmplDir = "template"

var (
	mdTmp  = filepath.Join(tmplDir, "top.template.md")
	orgTmp = filepath.Join(tmplDir, "org.template.org")
)
