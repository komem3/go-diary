package diary

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/manifoldco/promptui"
	"github.com/rakyll/statik/fs"
)

// Initialize create template files.
func Initialize(l Logger, dir string, tmplDir string, files []string) error {
	err := os.Mkdir(filepath.Join(dir, tmplDir), 0755)
	if err != nil {
		if os.IsExist(err) {
			yes, err := yesNoPrompt("Already exsits "+dir+". Do you overwrite this?", false)
			if err != nil {
				return fmt.Errorf("overwrite question: %w", err)
			}
			if !yes {
				return nil
			}
		} else {
			return fmt.Errorf("make template directory: %w", err)
		}
	}

	statikFS, err := fs.New()
	if err != nil {
		return fmt.Errorf("load static file data: %w", err)
	}

	for _, file := range files {
		l.Debug(
			"msg", "create "+file,
		)
		if err = createStaticFS(l, statikFS, filepath.Join(dir, file)); err != nil {
			return fmt.Errorf("create file %s: %w", file, err)
		}
	}
	return nil
}

func createStaticFS(l Logger, statikFS http.FileSystem, filPath string) error {
	r, err := statikFS.Open("/" + filepath.Base(filPath))
	if err != nil {
		return err
	}
	defer CloseWithErrLog(l, r)
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(filPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer CloseWithErrLog(l, file)

	_, err = file.Write(contents)
	return err
}

func yesNoPrompt(msg string, defaultYes bool) (bool, error) {
	var cursor int
	if !defaultYes {
		cursor = 1
	}
	prompt := promptui.Select{
		Label:     msg + "[Yes/No]",
		Items:     []string{"Yes", "No"},
		Stdin:     os.Stdin,
		CursorPos: cursor,
	}
	_, result, err := prompt.Run()
	return result == "Yes", err
}
