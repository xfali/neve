// Copyright (C) 2019-2020, Xiongfa Li.
// @author xiongfa.li
// @version V1.0
// Description:

package test

import (
	"github.com/gin-gonic/gin"
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
	neve.RegisterProcessor(ginImpl.NewProcessor(), processor.NewValueProcessor())

	app := neve.NewFileConfigApplication("assets/config-test.json")
	app.RegisterBean(&webBean{})
	app.Run()
}