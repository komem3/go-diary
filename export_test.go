package diary

import "time"

func (c *Creator) SetNowFunc(now func() time.Time) {
	c.now = now
}

func (g *Formatter) SetNowFunc(now func() time.Time) {
	g.now = now
}

var (
	FormatterCopyfFile = (*Formatter).copyFile
)
