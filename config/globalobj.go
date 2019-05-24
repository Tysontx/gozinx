package config

import (
	"io/ioutil"
	"encoding/json"
)

// 全局配置文件的类
type GlobalObj struct {
	Host            string // 当前监听的 IP
	Port            int    // 当前监听的 post
	Name            string // 当前 zinx server 名字
	Version         string // 当前框架版本号
	MaxPageageCount uint32 // 每次 read 一次的最大长度
}

// 定义一个全局对外的配置对象
var GlobalObject *GlobalObj

// 添加一个加载配置文件的方法
func (g *GlobalObj) LoadConfig() {
	data, err := ioutil.ReadFile("conf/gozinx.json")
	if err != nil {
		panic(err)
	}
	// 将 gozinx.json 的数据转换到 GlobalObj 中，json 解析
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

// 只要 import 当前模块，就会执行 init 函数，加载配置文件
func init() {
	// 配置文件的读取操作
	GlobalObject = &GlobalObj{
		Host:            "0.0.0.0",
		Port:            8999,
		Name:            "my gozinx app",
		Version:         "v0.4",
		MaxPageageCount: 512,
	}
	// 加载配置文件
	GlobalObject.LoadConfig()
}
