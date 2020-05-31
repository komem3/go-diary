package diary

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"
)

// Generator generate file
type Generator struct {
	logger Logger
	now    func() time.Time
	Err    error
}

// NewGenerator generate Generator
func NewGenerator(l Logger) *Generator {
	return &Generator{
		logger: l,
		now:    time.Now,
	}
}

// NewDiary generate today diary
func (g Generator) NewDiary(tmplFile, dir, nameFormat string) {
	if g.Err != nil {
		return
	}
	g.logger.Debug(
		"msg", "start to generate diary",
		"templateFile", tmplFile,
	)

	temp, err := template.New(filepath.Base(tmplFile)).ParseFiles(tmplFile)
	if err != nil {
		g.Err = err
		return
	}

	tmpName := fmt.Sprintf("diary-%v-tmp.txt", g.now().UnixNano())
	file, err := os.Create(tmpName)
	if err != nil {
		g.Err = err
		return
	}
	defer CloseWithErrLog(g.logger, file)

	now := g.now()
	err = temp.Execute(file, map[string]interface{}{
		"Year":    fmt.Sprintf("%d", now.Year()),
		"Month":   fmt.Sprintf("%02d", now.Month()),
		"Day":     fmt.Sprintf("%02d", now.Day()),
		"Weekday": now.Weekday().String(),
	})
	if err != nil {
		g.Err = err
		if err = os.Remove(tmpName); err != nil {
			g.logger.Error("err", err)
		}
		return
	}

	diaryFile := now.Format(nameFormat)
	err = os.Rename(tmpName, filepath.Join(dir, diaryFile))
	if err != nil {
		g.Err = err
		return
	}
	g.logger.Debug(
		"msg", "end to generate diary",
		"file", diaryFile,
	)
}
