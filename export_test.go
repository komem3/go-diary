package diary

import "time"

func (g *Formatter) SetNowFunc(now func() time.Time) {
	g.now = now
}

var (
	FormatterCopyfFile = (*Formatter).copyFile
)
