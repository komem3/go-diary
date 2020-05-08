package diary

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/golang/glog"
	"github.com/stretchr/testify/assert"
)

func TestWriteDirTree(t *testing.T) {
	type (
		args struct {
			elem         TopElem
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
				elem: TopElem{
					Years: []YearElem{
						{
							Year: "2019",
							Months: []MonthElem{
								{
									Month: "1",
									Days: []DayElem{
										{
											Day:  "01",
											Path: "2019/01/01",
										},
									},
								},
								{
									Month: "2",
									Days: []DayElem{
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
				elem:         TopElem{},
				readMePath:   "./testdata/WriteDirTree/top_readme_test2_in.txt",
				templatePath: "./testdata/WriteDirTree/top.template.md",
			},
			want{
				outputFile: "./testdata/WriteDirTree/top_readme_test2_out.txt",
				hasErr:     false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertions := assert.New(t)
			dGen := &DiayGeneretor{
				file:     tt.args.readMePath + ".copy",
				tmplFile: tt.args.templatePath,
				logger:   NewLogger(false),
				now:      time.Now,
				to:       ".",
			}
			err := dGen.copyFile(tt.args.readMePath, tt.args.readMePath+".copy")
			if !assertions.NoError(err) {
				return
			}
			dGen = dGen.WriteDirTree(tt.args.elem)
			if !assertions.NoError(dGen.Err) {
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

func TestParseFileMap(t *testing.T) {
	type (
		want struct {
			elem FileMap
		}
	)
	tests := []struct {
		name string
		args string
		want want
	}{
		{
			"success1",
			"./testdata/ParseDirTree/data1",
			want{
				FileMap{
					"2018": {
						"03": {
							"02": "testdata/ParseDirTree/data1/2018/3/20180302.md",
						},
					},
					"2019": {
						"01": {
							"01": "testdata/ParseDirTree/data1/2019/1/20190101.md",
							"02": "testdata/ParseDirTree/data1/2019/1/20190102_test.md"},
						"02": {
							"20": "testdata/ParseDirTree/data1/2019/2/20190220.md",
						},
					},
				},
			},
		},
		{
			"success2",
			"./testdata/ParseDirTree/data2",
			want{
				FileMap{
					"2020": {
						"02": {
							"03": "testdata/ParseDirTree/data2/20200203.md",
						},
						"04": {
							"04": "testdata/ParseDirTree/data2/2019/20200404.md",
						},
					},
					"2019": {
						"12": {
							"12": "testdata/ParseDirTree/data2/2019/20191212.md",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		assertions := assert.New(t)
		fmap := ParseFileMap(tt.args)
		assertions.Equal(tt.want.elem, fmap)
	}
}

func TestFormatDir(t *testing.T) {
	type args struct {
		fMap  FileMap
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
				FileMap{
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
				FileMap{
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
				FileMap{
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertions := assert.New(t)
			dGen := &DiayGeneretor{
				to:     tt.args.to,
				move:   tt.args.move,
				logger: NewLogger(false),
			}
			for src, dst := range tt.args.files {
				if err := dGen.copyFile(src, dst); !assertions.NoError(err) {
					return
				}
			}
			assertions.NoError(dGen.FormatDir(tt.args.fMap).Err)
			glog.Flush()
		})
	}
}

func TestMap2Elem(t *testing.T) {
	type (
		args struct {
			fMap FileMap
		}
		want struct {
			elem TopElem
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
				FileMap{
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
				TopElem{
					Years: []YearElem{
						{
							Year: "2019",
							Months: []MonthElem{
								{
									Month: "01",
									Days: []DayElem{
										{
											Day:  "20190101.md",
											Path: "2019/01/20190101.md",
										},
									},
								},
								{
									Month: "12",
									Days: []DayElem{
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
							Months: []MonthElem{
								{
									Month: "02",
									Days: []DayElem{
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertions := assert.New(t)
			assertions.Equal(tt.want.elem, Map2Elem(tt.args.fMap))
		})
	}
}
