package diary_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/komem3/go-diary"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFormatter_WriteDirTree(t *testing.T) {
	type (
		args struct {
			elem         diary.TopElem
			readMePath   string
			templatePath string
		}
		want struct {
			outputFile string
			hasErr     bool
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
				elem: diary.TopElem{
					Years: []diary.YearElem{
						{
							Year: "2019",
							Months: []diary.MonthElem{
								{
									Month: "1",
									Days: []diary.DayElem{
										{
											Day:  "01",
											Path: "2019/01/01",
										},
									},
								},
								{
									Month: "2",
									Days: []diary.DayElem{
										{
											Day:  "02",
											Path: "2019/02/02",
										},
									},
								},
							},
						},
						{
							Year: "2018",
						},
					},
				},
				readMePath:   "./testdata/WriteDirTree/top_readme_test1_in.txt",
				templatePath: "./testdata/WriteDirTree/top.template.md",
			},
			want{
				outputFile: "./testdata/WriteDirTree/top_readme_test1_out.txt",
				hasErr:     false,
			},
		},
		{
			"success2",
			args{
				elem:         diary.TopElem{},
				readMePath:   "./testdata/WriteDirTree/top_readme_test2_in.txt",
				templatePath: "./testdata/WriteDirTree/top.template.md",
			},
			want{
				outputFile: "./testdata/WriteDirTree/top_readme_test2_out.txt",
				hasErr:     false,
			},
		},
	}
	logger := diary.NewLogger(true)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertions := assert.New(t)

			f := diary.NewFormatter(logger)
			err := diary.FormatterCopyfFile(f, tt.args.readMePath, tt.args.readMePath+".copy")
			if !assertions.NoError(err) {
				return
			}
			f = f.WriteDirTree(
				tt.args.elem,
				tt.args.readMePath+".copy",
				tt.args.templatePath,
				".",
			)
			if !assertions.NoError(f.Err) {
				return
			}

			expectOut, err := ioutil.ReadFile(tt.want.outputFile)
			if !assertions.NoError(err) {
				return
			}
			output, err := ioutil.ReadFile(tt.args.readMePath + ".copy")
			if !assertions.NoError(err) {
				return
			}
			assertions.Equal(expectOut, output)
		})
	}
}

func TestFormatter_ParseFileMap(t *testing.T) {
	type (
		env struct {
			pwd string
		}
		want struct {
			elem diary.FileMap
		}
	)
	tests := []struct {
		name string
		env  env
		args string
		want want
	}{
		{
			"success",
			env{
				pwd: "",
			},
			"./testdata/ParseDirTree/data1",
			want{
				diary.FileMap{
					"2018": {
						"03": {
							"02": "testdata/ParseDirTree/data1/2018/3/20180302.md",
						},
					},
					"2019": {
						"01": {
							"01": "testdata/ParseDirTree/data1/2019/1/20190101.md",
							"02": "testdata/ParseDirTree/data1/2019/1/20190102_test.md",
						},
						"02": {
							"20": "testdata/ParseDirTree/data1/2019/2/20190220.md",
						},
					},
				},
			},
		},
		{
			"contain ignore file and dir",
			env{
				pwd: "testdata/ParseDirTree/data2",
			},
			".",
			want{
				diary.FileMap{
					"2020": {
						"02": {
							"03": "20200203.md",
						},
						"04": {
							"04": "2019/20200404.md",
						},
					},
					"2019": {
						"12": {
							"12": "2019/20191212.md",
						},
					},
				},
			},
		},
		{
			"parent directory",
			env{
				pwd: "testdata/ParseDirTree/data1/2018",
			},
			"..",
			want{
				diary.FileMap{
					"2018": {
						"03": {
							"02": "../2018/3/20180302.md",
						},
					},
					"2019": {
						"01": {
							"01": "../2019/1/20190101.md",
							"02": "../2019/1/20190102_test.md",
						},
						"02": {
							"20": "../2019/2/20190220.md",
						},
					},
				},
			},
		},
	}
	logger := diary.NewLogger(true)
	_, filename, _, _ := runtime.Caller(0)
	pwd := filepath.Dir(filename)
	defer func() { _ = os.Chdir(pwd) }()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			assertions := assert.New(t)
			err := os.Chdir(filepath.Join(pwd, tt.env.pwd))
			require.NoError(t, err)
			f := diary.NewFormatter(logger)
			fmap := f.ParseFileMap(tt.args)
			assertions.Equal(tt.want.elem, fmap)
		})
	}
}

func TestFormatter_FormatDir(t *testing.T) {
	type args struct {
		fMap  diary.FileMap
		files map[string]string
		to    string
		move  bool
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"sucess1",
			args{
				diary.FileMap{
					"2019": {
						"12": {
							"01": "testdata/FormatDir/input1/2019/12/20191201.md",
						},
					},
				},
				map[string]string{
					"testdata/FormatDir/input1/20191201.md": "testdata/FormatDir/input1/2019/12/20191201.md",
				},
				"testdata/FormatDir/output1",
				true,
			},
		},
		{
			"sucess2",
			args{
				diary.FileMap{
					"2018": {
						"01": {
							"01": "testdata/FormatDir/input2/2018/01/20180101.md",
							"02": "testdata/FormatDir/input2/2018/01/20180102.md",
						},
						"02": {
							"02": "testdata/FormatDir/input2/2018/02/20180202.md",
						},
					},
					"2019": {
						"12": {
							"12": "testdata/FormatDir/input2/2019/20191212.md",
						},
					},
				},
				map[string]string{
					"testdata/FormatDir/input2/20180101.md": "testdata/FormatDir/input2/2018/01/20180101.md",
					"testdata/FormatDir/input2/20180102.md": "testdata/FormatDir/input2/2018/01/20180102.md",
					"testdata/FormatDir/input2/20180202.md": "testdata/FormatDir/input2/2018/02/20180202.md",
					"testdata/FormatDir/input2/20191212.md": "testdata/FormatDir/input2/2019/20191212.md",
				},
				"testdata/FormatDir/output2",
				true,
			},
		},
		{
			"success_copy",
			args{
				diary.FileMap{
					"2018": {
						"01": {
							"01": "testdata/FormatDir/input2/2018/01/20180101.md",
							"02": "testdata/FormatDir/input2/2018/01/20180102.md",
						},
						"02": {
							"02": "testdata/FormatDir/input2/2018/02/20180202.md",
						},
					},
					"2019": {
						"12": {
							"12": "testdata/FormatDir/input2/2019/20191212.md",
						},
					},
				},
				map[string]string{
					"testdata/FormatDir/input2/20180101.md": "testdata/FormatDir/input2/2018/01/20180101.md",
					"testdata/FormatDir/input2/20180102.md": "testdata/FormatDir/input2/2018/01/20180102.md",
					"testdata/FormatDir/input2/20180202.md": "testdata/FormatDir/input2/2018/02/20180202.md",
					"testdata/FormatDir/input2/20191212.md": "testdata/FormatDir/input2/2019/20191212.md",
				},
				"testdata/FormatDir/output2",
				false,
			},
		},
	}
	logger := diary.NewLogger(true)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertions := assert.New(t)
			f := diary.NewFormatter(logger)
			for src, dst := range tt.args.files {
				if err := diary.FormatterCopyfFile(f, src, dst); !assertions.NoError(err) {
					return
				}
			}
			assertions.NoError(f.FormatDir(tt.args.fMap, tt.args.to, tt.args.move).Err)
		})
	}
}

func TestFormatter_Map2Elem(t *testing.T) {
	type (
		args struct {
			fMap    diary.FileMap
			options []diary.Map2ElemOptionFunc
		}
		want struct {
			elem diary.TopElem
		}
	)
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			"success_no_option",
			args{
				fMap: diary.FileMap{
					"2019": {
						"12": {
							"01": "test1/20191201.md",
							"12": "test1/20191212.md",
						},
						"01": {
							"01": "test1/20190101.md",
						},
					},
					"2018": {
						"02": {
							"31": "test1/20180231.md",
						},
					},
				},
			},
			want{
				diary.TopElem{
					Years: []diary.YearElem{
						{
							Year: "2019",
							Months: []diary.MonthElem{
								{
									Month: "01",
									Days: []diary.DayElem{
										{
											Day:  "20190101.md",
											Path: "2019/01/20190101.md",
										},
									},
								},
								{
									Month: "12",
									Days: []diary.DayElem{
										{
											Day:  "20191201.md",
											Path: "2019/12/20191201.md",
										},
										{
											Day:  "20191212.md",
											Path: "2019/12/20191212.md",
										},
									},
								},
							},
						},
						{
							Year: "2018",
							Months: []diary.MonthElem{
								{
									Month: "02",
									Days: []diary.DayElem{
										{
											Day:  "20180231.md",
											Path: "2018/02/20180231.md",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			"success_same_as_default_option",
			args{
				fMap: diary.FileMap{
					"2019": {
						"12": {
							"01": "test1/20191201.md",
							"12": "test1/20191212.md",
						},
						"01": {
							"01": "test1/20190101.md",
						},
					},
					"2018": {
						"02": {
							"31": "test1/20180231.md",
						},
					},
				},
				options: []diary.Map2ElemOptionFunc{
					diary.YearSort(diary.DESCSort),
					diary.MonthSort(diary.ASCSort),
					diary.DaySort(diary.ASCSort),
				},
			},
			want{
				diary.TopElem{
					Years: []diary.YearElem{
						{
							Year: "2019",
							Months: []diary.MonthElem{
								{
									Month: "01",
									Days: []diary.DayElem{
										{
											Day:  "20190101.md",
											Path: "2019/01/20190101.md",
										},
									},
								},
								{
									Month: "12",
									Days: []diary.DayElem{
										{
											Day:  "20191201.md",
											Path: "2019/12/20191201.md",
										},
										{
											Day:  "20191212.md",
											Path: "2019/12/20191212.md",
										},
									},
								},
							},
						},
						{
							Year: "2018",
							Months: []diary.MonthElem{
								{
									Month: "02",
									Days: []diary.DayElem{
										{
											Day:  "20180231.md",
											Path: "2018/02/20180231.md",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			"success_reverse_option",
			args{
				fMap: diary.FileMap{
					"2019": {
						"12": {
							"01": "test1/20191201.md",
							"12": "test1/20191212.md",
						},
						"01": {
							"01": "test1/20190101.md",
						},
					},
					"2018": {
						"02": {
							"31": "test1/20180231.md",
						},
					},
				},
				options: []diary.Map2ElemOptionFunc{
					diary.YearSort(diary.ASCSort),
					diary.MonthSort(diary.DESCSort),
					diary.DaySort(diary.DESCSort),
				},
			},
			want{
				diary.TopElem{
					Years: []diary.YearElem{
						{
							Year: "2018",
							Months: []diary.MonthElem{
								{
									Month: "02",
									Days: []diary.DayElem{
										{
											Day:  "20180231.md",
											Path: "2018/02/20180231.md",
										},
									},
								},
							},
						},
						{
							Year: "2019",
							Months: []diary.MonthElem{
								{
									Month: "12",
									Days: []diary.DayElem{
										{
											Day:  "20191212.md",
											Path: "2019/12/20191212.md",
										},
										{
											Day:  "20191201.md",
											Path: "2019/12/20191201.md",
										},
									},
								},
								{
									Month: "01",
									Days: []diary.DayElem{
										{
											Day:  "20190101.md",
											Path: "2019/01/20190101.md",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	logger := diary.NewLogger(true)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertions := assert.New(t)
			f := diary.NewFormatter(logger)
			assertions.Equal(tt.want.elem, f.Map2Elem(tt.args.fMap, tt.args.options...))
		})
	}
}
