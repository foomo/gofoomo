package bert

import (
	"github.com/foomo/gofoomo/foomo"
)

type Bert struct {
	foomo *foomo.Foomo
}

func NewBert(f *foomo.Foomo) *Bert {
	return &Bert{
		foomo: f,
	}
}

func (b *Bert) Hiccup() (report string, err error) {
	report = ""
	b.foomo.GetVarDir()
	return report, err
}
