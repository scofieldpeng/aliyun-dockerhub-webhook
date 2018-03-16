package main

// 镜像配置
type (
	Image struct {
		// app名称
		AppName string
		// 标签名
		Tag string
		// 镜像拉取地址
		ImageHost string
		// 脚本名称
		Script        string
		RemoteRequest struct {
			Url    string
			Method Method
		}
		Method string
	}
	Method string
)

const (
	Get    Method = "GET"
	Post   Method = "POST"
	Put    Method = "PUT"
	Delete Method = "DELETE"
)

// 执行脚本命令
func (c Image) runScript() (err error) {
	if c.Script == "" {
		return
	}
	
}
