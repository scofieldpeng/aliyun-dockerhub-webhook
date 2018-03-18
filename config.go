package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/vaughan0/go-ini"
)

// 配置文件
type config struct {
	appPath string
	// 配置路径
	configPath   string
	locker       *sync.RWMutex
	refreshToken string
	httpAddr     string
	// 配置
}

var (
	cfg = config{
		locker: &sync.RWMutex{},
	}
)

// 初始化config配置
func initConfig() (err error) {
	if Debug {
		fmt.Println("[app]读取配置文件")
	}
	cfg.appPath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	if cfg.configPath == "" {
		cfg.configPath = fmt.Sprintf("%s/config/app.ini", cfg.appPath)
	}
	if err = cfg.Reload(); err != nil {
		return
	}

	return
}

// 载入配置
func (c *config) Reload() (err error) {
	if c.configPath == "" {
		return errors.New("config path not set")
	}
	c.locker.Lock()
	defer c.locker.Unlock()
	var (
		data ini.File
	)
	data, err = ini.LoadFile(c.configPath)
	img := c.configToImage(data)
	images.Writes(img)

	c.refreshToken, _ = data.Get("system", "refresh_token")

	return
}

// 转化为image配置
func (c *config) configToImage(data ini.File) (images map[string]Image) {
	images = make(map[string]Image)
	for k, v := range data {
		if k != "system" {
			image := Image{
				AppName: k,
			}
			image.Debug = Debug
			image.AppToken = v["token"]
			image.ImageHost = v["host"]
			image.Tag = v["tag"]
			if image.Tag == "" {
				image.Tag = "all"
			}
			image.ScriptCallback = v["script"]
			if image.ScriptCallback == "" {
				continue
			}
			images[k] = image
		}
	}

	return
}
