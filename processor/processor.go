// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package processor

import (
	"github.com/xfali/fig"
	"io"
)

type Processor interface {
	Init(conf fig.Properties) error
	HandleBean(o interface{}) (bool, error)
	Process() error

	io.Closer
}
