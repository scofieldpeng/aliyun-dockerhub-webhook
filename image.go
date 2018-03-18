package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

// 镜像配置
type (
	Image struct {
		// app名称
		AppName string
		// 更新token
		AppToken string
		// 镜像名
		ImageName string
		// 标签名
		Tag string
		// 镜像拉取地址
		ImageHost string
		// 脚本名称
		ScriptCallback string
		// 命令执行超时时间
		Timeout time.Duration
		// 是否是debug模式
		Debug bool
	}
)

// 执行脚本命令
func (c Image) runScript() (err error) {
	var (
		cmd *exec.Cmd
	)
	if c.ScriptCallback == "" {
		return
	}

	// todo: 脚本运行只考虑在unix下运行，不考虑windows
	scriptPath := c.ScriptCallback
	if !filepath.IsAbs(c.ScriptCallback) {
		scriptPath = cfg.appPath + "/script/" + c.ScriptCallback
	}
	if c.Timeout > 0 {
		ctx, cancelFns := context.WithTimeout(context.Background(), c.Timeout)
		defer cancelFns()
		cmd = exec.CommandContext(ctx, "/bin/sh", "-c", scriptPath)
	} else {
		fmt.Println("not timeout")
		cmd = exec.Command("/bin/sh", "-c", scriptPath)
	}

	// 传入镜像参数到脚本的环境变量
	cmd.Env = append(cmd.Env, fmt.Sprintf("IMAGE_NAME=%s", c.ImageName))
	cmd.Env = append(cmd.Env, fmt.Sprintf("IMAGE_TAG=%s", c.Tag))
	cmd.Env = append(cmd.Env, fmt.Sprintf("IMAGE_HOST=%s", c.ImageHost))

	if c.Debug {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	err = cmd.Run()

	return
}

// TODO 执行远程推送
// 因为远程推送有很多特殊的地方，比如请求要携带的参数，头部信息等等，无法完全做到自定义，因此这里用自定义脚本来做吧
func (c Image) runRemote() (err error) {
	return
}

type Images struct {
	data   map[string]Image
	locker *sync.RWMutex
}

// 获取镜像配置
func (i *Images) Get(imageName string) (image Image, exist bool) {
	i.locker.RLock()
	defer i.locker.RUnlock()

	image, exist = i.data[imageName]
	return
}

// 写入镜像配置
func (i *Images) Write(data Image) {
	i.locker.Lock()
	defer i.locker.Unlock()

	i.data[data.AppName] = data
}

func (i *Images) Writes(data map[string]Image) {
	i.locker.Lock()
	defer i.locker.Unlock()
	i.data = data
}

var (
	images = Images{
		data:   make(map[string]Image),
		locker: &sync.RWMutex{},
	}
)
