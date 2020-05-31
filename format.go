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

// Formatter is generator diary
type Formatter struct {
	logger Logger
	Err    error
	now    func() time.Time
}

// NewFormatter generate Formatter
func NewFormatter(logger Logger) *Formatter {
	return &Formatter{
		logger: logger,
		now:    time.Now,
	}
}

// WriteDirTree write directory tree
func (f *Formatter) WriteDirTree(elem TopElem, filePath, templatePath, to string) *Formatter {
	if f.Err != nil {
		return f
	}

	f.logger.Debug(
		"msg", "start write dir tree",
		"file", filePath,
		"templateFile", templatePath,
	)

	tmpName := fmt.Sprintf("diary-%v-tmp.txt", f.now().UnixNano())
	err := func() error {
		tmpFile, err := os.Create(tmpName)
		if err != nil {
			return fmt.Errorf("open temp file: %w", err)
		}

		writer := bufio.NewWriter(tmpFile)
		defer CloseWithErrLog(f.logger, tmpFile)

		file, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
		if err != nil {
			if !os.IsNotExist(err) {
				return fmt.Errorf("open base file: %w", err)
			}
		} else {
			defer CloseWithErrLog(f.logger, file)

			reader := bufio.NewReader(file)
			for {
				line, err := reader.ReadString('\n')
				if err == io.EOF || string(line) == "# diary record\n" {
					break
				}
				if err != nil {
					return fmt.Errorf("read base file: %w", err)
				}
				_, err = writer.WriteString(line)
				if err != nil {
					return fmt.Errorf("write base content to temp file: %w", err)
				}
			}
		}

		temp, err := template.New(filepath.Base(templatePath)).ParseFiles(templatePath)
		if err != nil {
			return fmt.Errorf("open template file: %w", err)
		}
		elem.Base = to
		err = temp.Execute(writer, elem)
		if err != nil {
			return fmt.Errorf("parse template: %w", err)
		}
		err = writer.Flush()
		if err != nil {
			return fmt.Errorf("temp file write flush: %w", err)
		}
		return nil
	}()
	if err != nil {
		f.Err = err
		if err = os.Remove(tmpName); err != nil {
			f.logger.Error("err", err.Error())
		}
		return f
	}
	if err = os.Rename(tmpName, filePath); err != nil {
		f.Err = fmt.Errorf("rename temp file to %s: %w", filePath, err)
	}
	return f
}

// ParseFileMap analys dir and parse FileMap
func (f Formatter) ParseFileMap(root string) FileMap {
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
func (f *Formatter) FormatDir(fMap FileMap, to string, move bool) *Formatter {
	if f.Err != nil {
		return f
	}
	f.logger.Debug(
		"msg", "start format directory",
		"dir", to,
	)

	err := os.Mkdir(to, 0755)
	if err != nil && !os.IsExist(err) {
		f.Err = fmt.Errorf("create top directory: %w", err)
		return f
	}
	for year, yMap := range fMap {
		err := os.Mkdir(fmt.Sprintf("%s/%s", to, year), 0755)
		if err != nil && !os.IsExist(err) {
			f.Err = fmt.Errorf("create sub directory: %w", err)
			return f
		}
		for month, mMap := range yMap {
			err := os.Mkdir(fmt.Sprintf("%s/%s/%s", to, year, month), 0755)
			if err != nil && !os.IsExist(err) {
				f.Err = fmt.Errorf("create sub directory: %w", err)
				return f
			}
			for _, path := range mMap {
				dst := fmt.Sprintf("%s/%s/%s/%s", to, year, month, filepath.Base(path))
				if move {
					f.logger.Debug(
						"msg", "move file",
						"from", path,
						"to", dst,
					)
					err = os.Rename(path, dst)
					if err != nil {
						f.Err = fmt.Errorf("move from %s to %s: %w", path, dst, err)
					}
				} else {
					if path == dst {
						f.logger.Warn(
							"msg", "copy file is same to base file",
							"from", path,
							"to", dst,
						)
						continue
					}
					f.logger.Debug(
						"msg", "copy file",
						"from", path,
						"to", dst,
					)
					err = f.copyFile(path, dst)
					if err != nil {
						f.Err = fmt.Errorf("copy from %s to %s: %w", path, dst, err)
					}
				}
				if f.Err != nil {
					return f
				}
			}
		}
	}
	return f
}

// Map2Elem convert FileMap to TopElem
func (f Formatter) Map2Elem(fMap FileMap) (elem TopElem) {
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

func (f Formatter) copyFile(srcName, dstName string) error {
	src, err := os.Open(srcName)
	if err != nil {
		return err
	}
	defer CloseWithErrLog(f.logger, src)

	dst, err := os.Create(dstName)
	if err != nil {
		return err
	}
	defer CloseWithErrLog(f.logger, dst)

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
