// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package inject

import (
	"errors"
	"github.com/xfali/neve/utils"
	"reflect"
	"sync"
)

type Container struct {
	pool sync.Map
}

func (c *Container) Register(o interface{}) {
	c.pool.Store(utils.GetObjectName(o), o)
}

func (c *Container) Inject(o interface{}) error {
	v := reflect.ValueOf(o)
	if v.Kind() == reflect.Interface {
		return c.injectInterface(v)
	}
	t := v.Type()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	if t.Kind() == reflect.Struct {
		return c.injectStruct(v)
	}
	return errors.New("Type Not support. ")
}

func (c *Container) injectField(v reflect.Value, name string) error {
	return nil
}

func (c *Container) injectInterface(v reflect.Value) error {
	vt := v.Type()
	o, ok := c.pool.Load(utils.GetTypeName(vt))
	if ok {
		safeSet(v, reflect.ValueOf(o))
		return nil
	} else {
		c.pool.Range(func(key, value interface{}) bool {
			ot := reflect.TypeOf(value)
			if ot.Implements(ot) {
				safeSet(v, reflect.ValueOf(value))
				c.pool.Store(utils.GetTypeName(vt), value)
				return false
			}
			return true
		})
		return nil
	}
	return nil
}

func (c *Container) injectStruct(v reflect.Value) error {
	vt := v.Type()
	o, ok := c.pool.Load(utils.GetTypeName(vt))
	if ok {
		if v.CanSet() {
			safeSet(v, reflect.ValueOf(o))
			return nil
		}
	}
	return nil
}

func safeSet(dest, src reflect.Value) {
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

