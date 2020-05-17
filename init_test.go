package diary_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/komem3/diary"
	_ "github.com/komem3/diary/statik"
	"github.com/stretchr/testify/assert"
)

func TestInitialize(t *testing.T) {
	type (
		kind int
		args struct {
			input string
			dir   string
		}
		want struct {
			kind kind
			err  error
		}
	)
	const (
		overwrite = iota + 1
		notOverwrite
		new
	)

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"overwrite",
			args{
				input: "k\n",
				dir:   "testdata/Initialize/overwrite/",
			},
			want{
				kind: overwrite,
				err:  nil,
			},
		},
		{
			"not_overwrite",
			args{
				input: "\n",
				dir:   "testdata/Initialize/not_overwrite/",
			},
			want{
				kind: notOverwrite,
				err:  nil,
			},
		},
		{
			"new",
			args{
				input: "\n",
				dir:   "testdata/Initialize/new/",
			},
			want{
				kind: new,
				err:  nil,
			},
		},
	}
	tmplFile := "template"
	files := []string{"template/top.template.md", "template/org.template.org"}
	logger := diary.NewLogger(true)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertions := assert.New(t)

			err := stubIO(
				tt.args.input,
				func() error { return diary.Initialize(logger, tt.args.dir, tmplFile, files) },
			)
			if tt.want.err != nil {
				assertions.EqualError(err, tt.want.err.Error())
				return
			}
			if !assertions.NoError(err) {
				return
			}
			for _, f := range files {
				file := filepath.Join(tt.args.dir, f)
				switch tt.want.kind {
				case overwrite, new:
					if assertions.FileExists(file) {
						if err := os.Remove(file); err != nil {
							logger.Error("err", err)
						}
					}
				case notOverwrite:
					if !assertions.NoFileExists(file) {
						if err := os.Remove(file); err != nil {
							logger.Error("err", err)
						}
					}
				}
			}
			if tt.want.kind == new {
				if err := os.Remove(filepath.Join(tt.args.dir, tmplFile)); err != nil {
					logger.Error("err", err)
				}
			}
		})
	}
}

func stubIO(inStr string, fn func() error) error {
	r, w, err := os.Pipe()
	if err != nil {
		return err
	}
	defer w.Close()
	_, err = w.WriteString(inStr + "\n")
	if err != nil {
		return err
	}
	os.Stdin = r
	return fn()
}
