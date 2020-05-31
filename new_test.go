package diary_test

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/komem3/diary"
	"github.com/stretchr/testify/assert"
)

func TestCreator_NewDiary(t *testing.T) {
	type (
		fields struct {
			dir        string
			tmplFile   string
			nameFormat string
			now        time.Time
		}
		want struct {
			err         error
			outputPath  string
			correctFile string
		}
	)
	tests := []struct {
		name   string
		fields fields
		want   want
	}{
		{
			"success1",
			fields{
				dir:        "testdata/NewDiary/data1",
				tmplFile:   "testdata/NewDiary/data1/diary1.template.md",
				nameFormat: "20060102.md",
				now:        time.Date(2019, 1, 1, 0, 0, 0, 0, time.Local),
			},
			want{
				err:         nil,
				outputPath:  "testdata/NewDiary/data1/20190101.md",
				correctFile: "testdata/NewDiary/data1/correct.md",
			},
		},
		{
			"success2",
			fields{
				dir:        "testdata/NewDiary/data2",
				tmplFile:   "testdata/NewDiary/data2/diary2.template.md",
				nameFormat: "20060102_sample.md",
				now:        time.Date(2021, 12, 12, 12, 0, 0, 0, time.Local),
			},
			want{
				err:         nil,
				outputPath:  "testdata/NewDiary/data2/20211212_sample.md",
				correctFile: "testdata/NewDiary/data2/correct.md",
			},
		},
		{
			"error no tmplFile",
			fields{
				dir:        "testdata/NewDiary/error",
				tmplFile:   "tesdata/NewDiary/error/not_found.md",
				nameFormat: "20060102_simple.md",
				now:        time.Date(2021, 12, 12, 12, 0, 0, 0, time.Local),
			},
			want{
				err:         errors.New("open template file: open tesdata/NewDiary/error/not_found.md: no such file or directory"),
				outputPath:  "",
				correctFile: "",
			},
		},
		{
			"Ignore unregistered variables",
			fields{
				dir:        "testdata/NewDiary/data3",
				tmplFile:   "testdata/NewDiary/data3/ignore.template.md",
				nameFormat: "20060102.md",
				now:        time.Date(2018, 01, 12, 12, 0, 0, 0, time.Local),
			},
			want{
				err:         nil,
				outputPath:  "testdata/NewDiary/data3/20180112.md",
				correctFile: "testdata/NewDiary/data3/correct.md",
			},
		},
	}
	logger := diary.NewLogger(true)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertions := assert.New(t)
			creator := diary.NewCreator(logger)
			creator.SetNowFunc(func() time.Time { return tt.fields.now })

			creator.NewDiary(tt.fields.tmplFile, tt.fields.dir, tt.fields.nameFormat)
			if tt.want.err != nil {
				assertions.EqualError(creator.Err, tt.want.err.Error())
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
