// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package neve

import (
	"github.com/xfali/fig"
	"github.com/xfali/neve/container"
	"github.com/xfali/neve/ctx"
	"github.com/xfali/neve/log"
	"github.com/xfali/neve/processor"
	"github.com/xfali/xlog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Application interface {
	RegisterBean(o interface{}) error
	RegisterBeanByName(name string, o interface{}) error
	Run() error
}

type FileConfigApplication struct {
	config fig.Properties
	ctx    ctx.ApplicationContext
}

type Opt func(*FileConfigApplication)

func NewFileConfigApplication(configPath string, opts ...Opt) *FileConfigApplication {
	logger := log.GetLogger()
	prop, err := fig.LoadYamlFile(configPath)
	if err != nil {
		logger.Fatalln("load config file failed: ", err)
		return nil
	}
	c := container.New()
	ret := &FileConfigApplication{
		config: prop,
		ctx:    ctx.NewDefaultApplicationContext(ctx.OptSetContainer(c)),
	}

	for _, opt := range opts {
		opt(ret)
	}

	for _, v := range processors {
		v.Init(prop, c)
		ret.ctx.AddProcessor(v)
	}
	return ret
}

func (app *FileConfigApplication) RegisterBean(o interface{}) error {
	return app.ctx.RegisterBean(o)
}

func (app *FileConfigApplication) RegisterBeanByName(name string, o interface{}) error {
	return app.ctx.RegisterBeanByName(name, o)
}

func (app *FileConfigApplication) Run() error {
	app.ctx.NotifyListeners(ctx.ApplicationEventInitialized)
	HandlerSignal(app.ctx.Close)
	return nil
}

func OptSetApplicationContext(ctx ctx.ApplicationContext) Opt {
	return func(application *FileConfigApplication) {
		application.ctx = ctx
	}
}

func HandlerSignal(closers ...func() error) {
	var (
		ch = make(chan os.Signal, 1)
	)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		si := <-ch
		switch si {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			time.Sleep(time.Second * 2)
			xlog.Infof("get a signal %s, stop the server", si.String())
			for i := range closers {
				closers[i]()
			}
			time.Sleep(time.Second)
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
	xlog.Infof("------ Process exited ------")
}

var processors []processor.Processor

func RegisterProcessor(proc ...processor.Processor) {
	for _, v := range proc {
		if v != nil {
			processors = append(processors, v)
		}
	}
}
