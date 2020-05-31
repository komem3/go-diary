package diary

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"
)

// Creator generate file
type Creator struct {
	logger Logger
	now    func() time.Time
	Err    error
}

// NewCreator generate Creator
func NewCreator(l Logger) *Creator {
	return &Creator{
		logger: l,
		now:    time.Now,
	}
}

// NewDiary generate today diary
func (c *Creator) NewDiary(tmplFile, dir, nameFormat string) {
	if c.Err != nil {
		return
	}
	c.logger.Debug(
		"msg", "start to generate diary",
		"templateFile", tmplFile,
	)
	now := c.now()
	diaryFile := now.Format(nameFormat)
	if _, err := os.Stat(filepath.Join(dir, diaryFile)); !os.IsNotExist(err) {
		yes, err := yesNoPrompt("Already exsits "+diaryFile+". Do you overwrite this?", false)
		if err != nil {
			c.Err = fmt.Errorf("overwrite question: %w", err)
			return
		}
		if !yes {
			return
		}
	}

	temp, err := template.New(filepath.Base(tmplFile)).ParseFiles(tmplFile)
	if err != nil {
		c.Err = fmt.Errorf("open template file: %w", err)
		return
	}

	tmpName := fmt.Sprintf("diary-%v-tmp.txt", c.now().UnixNano())
	file, err := os.Create(tmpName)
	if err != nil {
		c.Err = fmt.Errorf("create diary file: %w", err)
		return
	}
	defer CloseWithErrLog(c.logger, file)

	err = temp.Execute(file, map[string]interface{}{
		"Year":    fmt.Sprintf("%d", now.Year()),
		"Month":   fmt.Sprintf("%02d", now.Month()),
		"Day":     fmt.Sprintf("%02d", now.Day()),
		"Weekday": now.Weekday().String(),
	})
	if err != nil {
		c.Err = fmt.Errorf("parse date information: %w", err)
		if err = os.Remove(tmpName); err != nil {
			c.logger.Error("err", err)
		}
		return
	}

	err = os.Rename(tmpName, filepath.Join(dir, diaryFile))
	if err != nil {
		c.Err = fmt.Errorf("move temp file to %s: %w", diaryFile, err)
		return
	}
	c.logger.Debug(
		"msg", "end to generate diary",
		"file", diaryFile,
	)
}
