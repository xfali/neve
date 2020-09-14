// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package utils

import (
	"reflect"
	"strings"
)

func GetObjectName(o interface{}) string {
	if o == nil {
		return ""
	}
	t := reflect.TypeOf(o)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.String {
		return o.(string)
	}

	return GetTypeName(t)
}

func GetTypeName(t reflect.Type) string {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	name := t.PkgPath()
	if name != "" {
		name = strings.Replace(name, "/", ".", -1) + "." + t.Name()
	}
	return name
}