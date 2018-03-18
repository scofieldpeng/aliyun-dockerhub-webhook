package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

var Debug bool

func init() {
	fmt.Println("[app]初始化appconfig")
	initAppConfig()
	if err := initConfig(); err != nil {
		fmt.Println("初始化config失败")
	}
}

func main() {
	initWebServer()
}

// 初始化appconfig配置
func initAppConfig() {
	flag.StringVar(&cfg.configPath, "config_path", "", "配置文件路径")
	flag.BoolVar(&Debug, "debug", false, "debug模式")
	flag.StringVar(&cfg.httpAddr, "http_addr", "0.0.0.0:2345", "监听的地址和端口")
	flag.Parse()

	// 支持从环境变量中读取配置
	if configPath := os.Getenv("ALI_WEBHOOK_CONFIG_PATH"); configPath != "" {
		cfg.configPath = configPath
	}
	if debug := os.Getenv("ALI_WEBHOOK_DEBUG"); debug != "" {
		if debugBool, err := strconv.ParseBool(debug); err == nil {
			Debug = debugBool
		}
	}
	if addr := os.Getenv("ALI_WEBHOOK_ADDR"); addr != "" {
		cfg.httpAddr = addr
	}
}
