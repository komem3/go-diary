package diary_test

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/komem3/go-diary"
	"github.com/stretchr/testify/assert"
)

func TestCreator_NewDiary(t *testing.T) {
	type (
		args struct {
			dir        string
			tmplFile   string
			nameFormat string
			now        time.Time
			yes        bool
		}
		want struct {
			errMsg      string
			outputPath  string
			correctFile string
		}
	)
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"success1",
			args{
				dir:        "testdata/NewDiary/data1",
				tmplFile:   "testdata/NewDiary/data1/diary1.template.md",
				nameFormat: "20060102.md",
				now:        time.Date(2019, 1, 1, 0, 0, 0, 0, time.Local),
				yes:        true,
			},
			want{
				outputPath:  "testdata/NewDiary/data1/20190101.md",
				correctFile: "testdata/NewDiary/data1/correct.md",
			},
		},
		{
			"success2",
			args{
				dir:        "testdata/NewDiary/data2",
				tmplFile:   "testdata/NewDiary/data2/diary2.template.md",
				nameFormat: "20060102_sample.md",
				now:        time.Date(2021, 12, 12, 12, 0, 0, 0, time.Local),
				yes:        true,
			},
			want{
				outputPath:  "testdata/NewDiary/data2/20211212_sample.md",
				correctFile: "testdata/NewDiary/data2/correct.md",
			},
		},
		{
			"error no tmplFile",
			args{
				dir:        "testdata/NewDiary/error",
				tmplFile:   "tesdata/NewDiary/error/not_found.md",
				nameFormat: "20060102_simple.md",
				now:        time.Date(2021, 12, 12, 12, 0, 0, 0, time.Local),
				yes:        true,
			},
			want{
				errMsg:      "open template file: open tesdata/NewDiary/error/not_found.md: no such file or directory",
				outputPath:  "",
				correctFile: "",
			},
		},
		{
			"Ignore unregistered variables",
			args{
				dir:        "testdata/NewDiary/data3",
				tmplFile:   "testdata/NewDiary/data3/ignore.template.md",
				nameFormat: "20060102.md",
				now:        time.Date(2018, 01, 12, 12, 0, 0, 0, time.Local),
				yes:        true,
			},
			want{
				outputPath:  "testdata/NewDiary/data3/20180112.md",
				correctFile: "testdata/NewDiary/data3/correct.md",
			},
		},
		{
			"not overwrite",
			args{
				dir:        "testdata/NewDiary/data4",
				tmplFile:   "testdata/NewDiary/data4/overwrite.template.md",
				nameFormat: "20060102.md",
				now:        time.Date(2018, 01, 12, 12, 0, 0, 0, time.Local),
				yes:        false,
			},
			want{
				outputPath:  "testdata/NewDiary/data4/20180112.md",
				correctFile: "testdata/NewDiary/data4/correct.md",
			},
		},
	}
	logger := diary.NewLogger(true)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertions := assert.New(t)
			creator := diary.NewCreator(logger)
			creator.SetNowFunc(func() time.Time { return tt.args.now })

			err := stubIO(yesnoStr(tt.args.yes), func() error {
				path := creator.NewDiary(tt.args.tmplFile, tt.args.dir, tt.args.nameFormat)
				if tt.args.yes {
					assertions.Equal(tt.want.outputPath, path)
				} else {
					assertions.Empty(path)
				}
				return nil
			})
			if err != nil {
				panic(err)
			}
			if tt.want.errMsg != "" {
				assertions.EqualError(creator.Err, tt.want.errMsg)
				return
			}
			if !assertions.NoError(creator.Err) {
				return
			}

			output, err := ioutil.ReadFile(tt.want.outputPath)
			if !assertions.NoError(err) {
				if !os.IsNotExist(err) {
					err = os.Remove(tt.want.outputPath)
					assertions.NoError(err)
				}
				return
			}
			correct, err := ioutil.ReadFile(tt.want.correctFile)
			if !assertions.NoError(err) {
				return
			}
			assertions.Equal(string(correct), string(output))
		})
	}
}

func yesnoStr(yes bool) string {
	if yes {
		return "k\n"
	}
	return "\n"
}
