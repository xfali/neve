// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package injector

import (
	"errors"
	"github.com/xfali/neve/container"
	"github.com/xfali/neve/utils"
	"github.com/xfali/xlog"
	"reflect"
)

const (
	injectTagName = "inject"
)

type Injector interface {
	Inject(container container.Container, o interface{}) error
}

type DefaultInjector struct {
	Logger xlog.Logger
}

func New(log xlog.Logger) *DefaultInjector {
	return &DefaultInjector{
		Logger: log,
	}
}

func (injector *DefaultInjector) Inject(c container.Container, o interface{}) error {
	v := reflect.ValueOf(o)
	if v.Kind() == reflect.Interface {
		return injector.injectInterface(c, v, "")
	}
	t := v.Type()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		v = v.Elem()
	}

	if t.Kind() == reflect.Struct {
		return injector.injectStructFields(c, v)
	}
	return errors.New("Type Not support. ")
}

func (injector *DefaultInjector) injectStructFields(c container.Container, v reflect.Value) error {
	t := v.Type()
	if t.Kind() != reflect.Struct {
		return errors.New("result must be struct ptr")
	}

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		tag, ok := field.Tag.Lookup(injectTagName)
		if ok {
			fieldValue := v.Field(i)
			if fieldValue.Kind() == reflect.Ptr {
				fieldValue = fieldValue.Elem()
			}
			if fieldValue.CanSet() {
				switch fieldValue.Kind() {
				case reflect.Interface:
					err := injector.injectInterface(c, fieldValue, tag)
					if err != nil {
						injector.Logger.Errorln(err)
					}
				case reflect.Struct:
					err := injector.injectStruct(c, fieldValue, tag)
					if err != nil {
						injector.Logger.Errorln(err)
					}
				}
			}
		}
	}

	return nil
}

func (injector *DefaultInjector) injectInterface(c container.Container, v reflect.Value, name string) error {
	vt := v.Type()
	if name == "" {
		name = utils.GetTypeName(vt)
	}
	o, ok := c.Get(name)
	if ok {
		v.Set(reflect.ValueOf(o))
		return nil
	} else {
		c.Scan(func(key string, value interface{}) bool {
			ot := reflect.TypeOf(value)
			if ot.Implements(vt) {
				v.Set(reflect.ValueOf(value))
				c.RegisterByName(utils.GetTypeName(vt), value)
				return false
			}
			return true
		})
		return nil
	}
	return nil
}

func (injector *DefaultInjector) injectStruct(c container.Container, v reflect.Value, name string) error {
	vt := v.Type()
	if name == "" {
		name = utils.GetTypeName(vt)
	}
	o, ok := c.Get(name)
	if ok {
		utils.SafeSet(v, reflect.ValueOf(o))
		return nil
	}
	return nil
}