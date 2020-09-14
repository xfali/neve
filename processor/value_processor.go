// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package processor

import "github.com/xfali/fig"

type ValueProcessor struct {
	conf fig.Properties
}

func NewValueProcessor() *ValueProcessor {
	return &ValueProcessor{}
}

func (p *ValueProcessor) Init(conf fig.Properties) error {
	p.conf = conf
	return nil
}

func (p *ValueProcessor) HandleBean(o interface{}) (bool, error) {
	return true, fig.Fill(p.conf, o)
}

func (p *ValueProcessor) Process() error {
	return nil
}

func (p *ValueProcessor) Close() error {
	return nil
}
