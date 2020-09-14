// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package utils

import (
	"reflect"
	"strings"
)

func GetObjectName(o ...interface{}) string {
	if len(o) == 0 {
		return ""
	} else {
		if o[0] == nil {
			return ""
		}
		t := reflect.TypeOf(o[0])
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if t.Kind() == reflect.String {
			return o[0].(string)
		}

		name := t.PkgPath()
		if name != "" {
			name = strings.Replace(name, "/", ".", -1) + "." + t.Name()
		}
		return name
	}
}
