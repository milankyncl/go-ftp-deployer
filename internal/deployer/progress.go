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

func (p *Progress) Set(colors ...color.Attribute) {
	p.color = p.color.Add(colors...)
}

func (p *Progress) Message(msg ...interface{}) {
	_, _ = p.color.Println(msg...)
}
