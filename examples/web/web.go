// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/xfali/neve"
	"github.com/xfali/neve/log"
	"github.com/xfali/neve/processor"
	"github.com/xfali/neve/web/ginImpl"
	"github.com/xfali/neve/web/ginImpl/midware"
	"github.com/xfali/neve/web/result"
	"net/http"
)

func main() {
	neve.SetResourceRoot("examples/web")
	neve.RegisterProcessor(ginImpl.NewProcessor(), processor.NewValueProcessor())

	app := neve.NewFileConfigApplication(neve.GetResource("assets/config-test.json"))
	app.RegisterBean(&webServiceImpl{})
	app.RegisterBean(&webBean{})
	app.Run()
}

type webService interface {
	GetValue() string
}

type webServiceImpl struct {
	V string `fig:"userData.value"`
}

func (w *webServiceImpl) GetValue() string {
	return w.V
}

type webBean struct {
	Service webService `inject:""`
}

func (b *webBean) Register(engine gin.IRouter) {
	loghttp := midware.LogHttpUtil{
		Logger:      log.GetLogger(),
		LogRespBody: true,
	}
	engine.GET("test", loghttp.LogHttp(), func(context *gin.Context) {
		context.JSON(http.StatusOK, result.Ok(b.Service.GetValue()))
	})

	engine.GET("panic", loghttp.LogHttp(), func(context *gin.Context) {
		panic("test!")
	})
}
