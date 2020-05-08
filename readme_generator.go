package diary

import (
	"bufio"
	"errors"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"
)

var (
	ErrNotParamater = errors.New("insufficient parameters")
)

type FileMap map[Year]map[Month]map[Day]string

// DiaryGenerator is generator diary
type DiayGeneretor struct {
	from     string
	to       string
	file     string
	tmplFile string
	org      bool
	move     bool
	logger   Logger
	Err      error
	now      func() time.Time
}

// NewDiaryGenerator generate DiayGeneretor
func NewDiaryGenerator(file, tmplFile, from, to string, org bool, logger Logger) (*DiayGeneretor, error) {
	switch "" {
	case file, tmplFile, from:
		return nil, ErrNotParamater
	}

	if to == "" {
		to = from
	}
	return &DiayGeneretor{
		from:     from,
		to:       to,
		file:     file,
		tmplFile: tmplFile,
		org:      org,
		move:     from == to,
		logger:   logger,
		now:      time.Now,
	}, nil
}

// WriteDirTree write directory tree
func (d *DiayGeneretor) WriteDirTree(elem TopElem) *DiayGeneretor {
	if d.Err != nil {
		return d
	}

	d.logger.Debug(
		"msg", "start write dir tree",
		"file", d.file,
		"templateFile", d.tmplFile,
	)

	tmpName := fmt.Sprintf("diary-%v-tmp.txt", d.now().UnixNano())
	err := func() error {
		tmpFile, err := os.OpenFile(tmpName, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}

		writer := bufio.NewWriter(tmpFile)
		defer CloseWithErrLog(d.logger, tmpFile)

		file, err := os.OpenFile(d.file, os.O_RDONLY, 0644)
		if err != nil {
			if !os.IsNotExist(err) {
				return err
			}
		} else {
			defer CloseWithErrLog(d.logger, file)

			reader := bufio.NewReader(file)
			for {
				line, err := reader.ReadString('\n')
				if err == io.EOF ||
					(!d.org && string(line) == "# diary record\n") ||
					(d.org && string(line) == "* diary record\n") {
					break
				}
				if err != nil {
					return err
				}
				_, err = writer.WriteString(line)
				if err != nil {
					return err
				}
			}
		}

		temp, err := template.New(filepath.Base(d.tmplFile)).ParseFiles(d.tmplFile)
		if err != nil {
			return err
		}
		elem.Base = d.to
		err = temp.Execute(writer, elem)
		if err != nil {
			return err
		}
		err = writer.Flush()
		if err != nil {
			return err
		}
		return nil
	}()
	if err != nil {
		d.Err = err
		if err = os.Remove(tmpName); err != nil {
			d.logger.Error("err", err.Error())
		}
		return d
	}
	d.Err = os.Rename(tmpName, d.file)
	return d
}

// ParseFileMap analys dir and parse FileMap
func ParseFileMap(root string) FileMap {
	re := regexp.MustCompile(`([0-9]{4})([0-9]{2})([0-9]{2}).*\.[a-zA-Z]+$`)
	fmap := make(FileMap)
	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		submathes := re.FindStringSubmatch(path)
		if len(submathes) == 0 {
			return nil
		}
		year, month, day := Year(submathes[1]), Month(submathes[2]), Day(submathes[3])
		yMap, found := fmap[year]
		if !found {
			fmap[year] = map[Month]map[Day]string{}
			yMap = fmap[year]
		}
		mMap, found := yMap[month]
		if !found {
			yMap[month] = map[Day]string{}
			mMap = yMap[month]
		}
		mMap[day] = path
		return nil
	}); err != nil {
		panic(err)
	}
	return fmap
}

// FormatDir format directory
func (d *DiayGeneretor) FormatDir(fMap FileMap) *DiayGeneretor {
	if d.Err != nil {
		return d
	}
	d.logger.Debug(
		"msg", "start format directory",
		"dir", d.to,
	)

	err := os.Mkdir(d.to, 0755)
	if err != nil && !os.IsExist(err) {
		d.Err = err
		return d
	}
	for year, yMap := range fMap {
		err := os.Mkdir(fmt.Sprintf("%s/%s", d.to, year), 0755)
		if err != nil && !os.IsExist(err) {
			d.Err = err
			return d
		}
		for month, mMap := range yMap {
			err := os.Mkdir(fmt.Sprintf("%s/%s/%s", d.to, year, month), 0755)
			if err != nil && !os.IsExist(err) {
				d.Err = err
				return d
			}
			for _, path := range mMap {
				dst := fmt.Sprintf("%s/%s/%s/%s", d.to, year, month, filepath.Base(path))
				if d.move {
					d.logger.Debug(
						"msg", "move file",
						"from", path,
						"to", dst,
					)
					d.Err = os.Rename(path, dst)
				} else {
					if path == dst {
						d.logger.Warn(
							"msg", "copy file is same to base file",
							"from", path,
							"to", dst,
						)
						continue
					}
					d.logger.Debug(
						"msg", "copy file",
						"from", path,
						"to", dst,
					)
					d.Err = d.copyFile(path, dst)
				}
				if d.Err != nil {
					return d
				}
			}
		}
	}
	return d
}

// Map2Elem convert FileMap to TopElem
func Map2Elem(fMap FileMap) (elem TopElem) {
	var i, j int
	for y, yMap := range fMap {
		elem.Years = append(elem.Years, YearElem{Year: y})
		for m, mMap := range yMap {
			elem.Years[i].Months = append(elem.Years[i].Months, MonthElem{Month: m})
			for _, path := range mMap {
				elem.Years[i].Months[j].Days = append(
					elem.Years[i].Months[j].Days,
					DayElem{
						Day:  Day(filepath.Base(path)),
						Path: filepath.Join(string(y), string(m), filepath.Base(path)),
					},
				)
			}
			sort.Slice(elem.Years[i].Months[j].Days, func(l, r int) bool {
				return elem.Years[i].Months[j].Days[l].Day < elem.Years[i].Months[j].Days[r].Day
			})
			j++
		}
		j = 0
		sort.Slice(elem.Years[i].Months, func(l, r int) bool {
			return elem.Years[i].Months[l].Month < elem.Years[i].Months[r].Month
		})
		i++
	}
	sort.Slice(elem.Years, func(l, r int) bool {
		return elem.Years[l].Year > elem.Years[r].Year
	})

	return elem
}

func (d DiayGeneretor) copyFile(srcName, dstName string) error {
	src, err := os.Open(srcName)
	if err != nil {
		return err
	}
	defer CloseWithErrLog(d.logger, src)

	dst, err := os.Create(dstName)
	if err != nil {
		return err
	}
	defer CloseWithErrLog(d.logger, dst)

	_, err = io.Copy(dst, src)
	return err
}

type Closer interface {
	Close() error
}

func CloseWithErrLog(l Logger, c Closer) {
	err := c.Close()
	if err != nil {
		l.Error("err", err.Error())
	}
}
