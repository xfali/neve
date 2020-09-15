// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package container

import (
	"errors"
	"github.com/xfali/neve/utils"
	"reflect"
	"sync"
)

type Container interface {
	Register(o interface{}) error
	RegisterByName(name string, o interface{}) error

	Get(name string) (interface{}, bool)
	GetByType(o interface{}) bool

	Scan(func(key string, value interface{}) bool)
}

type DefaultContainer struct {
	objectPool sync.Map
}

func New() *DefaultContainer {
	return &DefaultContainer{}
}

func (c *DefaultContainer) Register(o interface{}) error {
	return c.RegisterByName(utils.GetObjectName(o), o)
}

func (c *DefaultContainer) RegisterByName(name string, o interface{}) error {
	t := reflect.TypeOf(o)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
		if t.Kind() == reflect.Ptr {
			return errors.New("o must be a Pointer but get Pointer's Pointer")
		}
	} else {
		return errors.New("o must be a Pointer")
	}

	c.objectPool.Store(name, o)
	return nil
}

func (c *DefaultContainer) Get(name string) (interface{}, bool) {
	return c.objectPool.Load(name)
}

func (c *DefaultContainer) GetByType(o interface{}) bool {
	v := reflect.ValueOf(o)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	o, ok := c.Get(utils.GetTypeName(v.Type()))
	if ok {
		utils.SafeSet(v, reflect.ValueOf(o))
	}
	return false
}

func (c *DefaultContainer) Scan(f func(key string, value interface{}) bool) {
	c.objectPool.Range(func(key, value interface{}) bool {
		return f(key.(string), value)
	})
}
