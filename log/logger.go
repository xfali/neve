// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package log

import "github.com/xfali/xlog"

type facFunc func(o ...interface{}) xlog.Logger

var factory = xlog.GetLogger

func GetLogger(o ...interface{}) xlog.Logger {
	return factory(o...)
}

func InitFactory(f facFunc) {
	factory = f
}
