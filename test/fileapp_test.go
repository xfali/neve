// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/xfali/fig"
	"github.com/xfali/neve"
	"github.com/xfali/neve/log"
	"github.com/xfali/neve/processor"
	"github.com/xfali/neve/web/ginImpl"
	"github.com/xfali/neve/web/ginImpl/midware"
	"github.com/xfali/neve/web/result"
	"net/http"
	"testing"
)

type webBean struct {
	V string `fig:"Log.Level"`
}

func (b *webBean) Print(str string) {
	fmt.Println(str)
}

func (b *webBean) Register(engine gin.IRouter) {
	loghttp := midware.LogHttpUtil{
		Logger:      log.GetLogger(),
		LogRespBody: true,
	}
	engine.GET("test", loghttp.LogHttp(), func(context *gin.Context) {
		context.JSON(http.StatusOK, result.Ok(b.V))
	})

	engine.GET("panic", loghttp.LogHttp(), func(context *gin.Context) {
		panic("test!")
	})
}

func TestWebAndValue(t *testing.T) {
	neve.RegisterProcessor(ginImpl.NewProcessor(), processor.NewValueProcessor(), &testProcess{})

	app := neve.NewFileConfigApplication("assets/config-test.json")
	app.RegisterBean(&webBean{})
	app.Run()
}

type testProcess struct{}

type print interface {
	Print(str string)
}

func (p *testProcess) Init(conf fig.Properties) error {
	return nil
}

func (p *testProcess) HandleBean(o interface{}) (bool, error) {
	switch v := o.(type) {
	case print:
		v.Print("test")
	}
	return true, nil
}

func (p *testProcess) Process() error {
	return nil
}

func (p *testProcess) Close() error {
	return nil
}
