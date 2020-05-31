package diary

import "time"

func (g *Generator) SetNowFunc(now func() time.Time) {
	g.now = now
}
