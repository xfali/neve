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
	Logger    xlog.Logger
	recursive bool
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
			fieldType := fieldValue.Type()
			if fieldType.Kind() == reflect.Ptr {
				fieldType = fieldType.Elem()
			}
			if fieldValue.CanSet() {
				switch fieldType.Kind() {
				case reflect.Interface:
					err := injector.injectInterface(c, fieldValue, tag)
					if err != nil {
						injector.Logger.Errorf("Inject Field error: [%s: %s] %s\n ", utils.GetTypeName(t), field.Name, err.Error())
					}
				case reflect.Struct:
					err := injector.injectStruct(c, fieldValue, tag)
					if err != nil {
						injector.Logger.Errorf("Inject Field error: [%s: %s] %s\n ", utils.GetTypeName(t), field.Name, err.Error())
					}
				}
			} else {
				injector.Logger.Errorf("Inject failed: Field cannot SET [%s: %s]\n ", utils.GetTypeName(t), field.Name)
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
		//自动注入
		var matchValues []interface{}
		c.Scan(func(key string, value interface{}) bool {
			//指定名称注册的对象直接跳过，因为在container.Get未满足，所以认定不是用户想要注入的对象
			if key != utils.GetObjectName(value) {
				return true
			}
			ot := reflect.TypeOf(value)
			if ot.Implements(vt) {
				matchValues = append(matchValues, value)
				if len(matchValues) > 1 {
					panic("Auto Inject bean found more than 1")
				}
				return true
			}
			return true
		})
		if len(matchValues) == 1 {
			v.Set(reflect.ValueOf(matchValues[0]))
			// cache to container
			err := c.RegisterByName(utils.GetTypeName(vt), matchValues[0])
			if err != nil {
				injector.Logger.Warnln(err)
			}
			return nil
		}
	}
	return errors.New("Inject nothing, cannot find any Implementation: " + utils.GetTypeName(vt))
}

func (injector *DefaultInjector) injectStruct(c container.Container, v reflect.Value, name string) error {
	vt := v.Type()
	if name == "" {
		name = utils.GetTypeName(vt)
	}
	o, ok := c.Get(name)
	if ok {
		ov := reflect.ValueOf(o)
		if vt.Kind() == reflect.Ptr {
			v.Set(ov)
		} else {
			v.Set(ov.Elem())
		}
		return nil
	}

	if injector.recursive {
		return injector.injectStructFields(c, v)
	} else {
		return errors.New("Inject nothing, cannot find any instance of  " + utils.GetTypeName(vt))
	}
}
