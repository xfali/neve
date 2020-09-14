// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package ctx

import (
	"github.com/xfali/neve/log"
	"github.com/xfali/neve/processor"
	"github.com/xfali/neve/utils"
	"github.com/xfali/xlog"
	"sync"
)

type ApplicationEvent int

const (
	ApplicationEventNone ApplicationEvent = iota
	ApplicationEventInitialized
)

type ApplicationContext interface {
	AddProcessor(processor.Processor)

	RegisterBean(o interface{}) error
	RegisterBeanByName(name string, o interface{}) error

	GetBean(name string) interface{}
	GetBeanByType(o interface{}) interface{}

	AddListener(ApplicationContextListener)

	NotifyListeners(ApplicationEvent)

	Close() error
}

type ApplicationContextListener interface {
	OnRefresh(ctx ApplicationContext)
	OnEvent(e ApplicationEvent, ctx ApplicationContext)
}

type DefaultApplicationContext struct {
	logger     xlog.Logger
	objectPool sync.Map

	listeners    []ApplicationContextListener
	listenerLock sync.Mutex

	processors []processor.Processor
}

func NewDefaultApplicationContext() *DefaultApplicationContext {
	return &DefaultApplicationContext{
		logger: log.GetLogger(),
	}
}

func (ctx *DefaultApplicationContext) Close() (err error) {
	for _, processor := range ctx.processors {
		err = processor.Close()
		if err != nil {
			ctx.logger.Errorln(err)
		}
	}
	return
}

func (ctx *DefaultApplicationContext) RegisterBean(o interface{}) error {
	return ctx.RegisterBeanByName(utils.GetObjectName(o), o)
}

func (ctx *DefaultApplicationContext) RegisterBeanByName(name string, o interface{}) error {
	ctx.objectPool.Store(name, o)

	for _, processor := range ctx.processors {
		_, err := processor.HandleBean(o)
		if err != nil {
			ctx.logger.Errorln(err)
		}
	}

	ctx.listenerLock.Lock()
	defer ctx.listenerLock.Unlock()

	for _, v := range ctx.listeners {
		v.OnRefresh(ctx)
	}

	switch o.(type) {
	case ApplicationContextListener:
		ctx.listeners = append(ctx.listeners, o.(ApplicationContextListener))
	}
	return nil
}

func (ctx *DefaultApplicationContext) GetBean(name string) interface{} {
	if o, ok := ctx.objectPool.Load(name); ok {
		return o
	}
	return nil
}

func (ctx *DefaultApplicationContext) AddProcessor(p processor.Processor) {
	ctx.processors = append(ctx.processors, p)
}

func (ctx *DefaultApplicationContext) AddListener(l ApplicationContextListener) {
	ctx.listenerLock.Lock()
	defer ctx.listenerLock.Unlock()

	ctx.listeners = append(ctx.listeners, l)
}

func (ctx *DefaultApplicationContext) NotifyListeners(e ApplicationEvent) {
	if ApplicationEventInitialized == e {
		for _, processor := range ctx.processors {
			processor.Process()
		}
	}

	ctx.listenerLock.Lock()
	defer ctx.listenerLock.Unlock()
	for _, v := range ctx.listeners {
		v.OnEvent(e, ctx)
	}
}

func (ctx *DefaultApplicationContext) GetBeanByType(o interface{}) interface{} {
	return ctx.GetBean(utils.GetObjectName(o))
}
