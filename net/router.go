package net

import "gozinx/ziface"

type BaseRouter struct {
}

// 在处理业务之前的方法
func (br *BaseRouter)PreHandle(request ziface.IRequest) {
	// 将 interface 的方法全部实现，目的是让用户重写这个方法
}
// 处理业务的主方法
func (br *BaseRouter)Handle(request ziface.IRequest) {
	// 将 interface 的方法全部实现，目的是让用户重写这个方法
}
// 处理业务之后的方法
func (br *BaseRouter)PostHandle(request ziface.IRequest) {
	// 将 interface 的方法全部实现，目的是让用户重写这个方法
}