package main

import (
	"fmt"
	"net/http"

	"gopkg.in/gin-gonic/gin.v1"
)

var (
	g = gin.New()
)

func initWebServer() (err error) {
	if Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	g.POST(`/webhook/:appname`, Webhook)
	g.GET(`/sys/config/refresh`, RefreshConfig)
	return http.ListenAndServe(cfg.httpAddr, g)
}

// webhook触发
func Webhook(ctx *gin.Context) {
	appName := ctx.Param("appname")
	if appName == "" {
		ctx.String(http.StatusBadRequest, "invalid app name")
		return
	}
	image, ok := images.Get(appName)
	if !ok {
		ctx.String(http.StatusNotFound, "app(%s) not found", appName)
		return
	}
	if image.AppToken != ctx.Query("token") {
		ctx.String(http.StatusForbidden, "forbidden visit")
		return
	}
	w := webhook{}
	if err := w.Parse(ctx); err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	image.ImageName = w.Image
	if image.Tag == "all" || image.Tag == "" || image.Tag == w.Tag {
		if image.ScriptCallback != "" {
			// 触发吧少年
			go func() {
				fmt.Println("trigger script,appname:", image.AppName)
				if err := image.runScript(); err != nil {
					fmt.Println(err)
				}
			}()
		}
	}

	ctx.String(http.StatusOK, "ok")
	return
}

// 刷新config
func RefreshConfig(ctx *gin.Context) {
	token := ctx.Query("token")
	if token == cfg.refreshToken {
		if err := cfg.Reload(); err != nil {
			ctx.String(http.StatusBadRequest, err.Error())
			return
		}
		ctx.String(http.StatusOK, "reload config success")
		return
	}

	ctx.String(http.StatusForbidden, "forbidden visit")
}
