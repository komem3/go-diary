package cmd

import "path/filepath"

const tmplDir = "template"

func mdTmp() string {
	return filepath.Join(tmplDir, "top.template.md")
}

func diaryTmp() string {
	return filepath.Join(tmplDir, "diary.template.md")
}
