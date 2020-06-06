package deployer

import (
	"github.com/fatih/color"
)

func NewProgress() *Progress {
	return &Progress{
		color: color.New(),
	}
}

type Progress struct {
	color *color.Color
}

func (p *Progress) Color(colors ...color.Attribute) *Progress {
	p.color = p.color.Add(colors...)
	return p
}

func (p *Progress) Reset() *Progress {
	p.color = p.color.Add(color.Reset)
	return p
}

func (p *Progress) Message(msg ...interface{}) *Progress {
	_, _ = p.color.Println(msg...)
	return p
}
