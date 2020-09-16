// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package refection

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

func SafeSet(dest, src reflect.Value) {
	if dest.Kind() == reflect.Ptr {
		dest = dest.Elem()
	}

	if src.Kind() == reflect.Ptr {
		src = src.Elem()
	}

	if dest.CanSet() {
		dest.Set(src)
	}
}
