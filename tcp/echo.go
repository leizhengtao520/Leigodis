package tcp

import (
	"context"
	"net"
	"sync"
	"Leigodis/atomic"
)

//客户端的链接抽象
type Client struct {
	//tcp 链接
	Conn net.Conn
	Watting Wait

}
type EchoHandler struct {
	// 保存所有工作状态client的集合(把map当set用)
	// 需使用并发安全的容器
	active sync.Map
	closing atomic.Boolean
}
func MakeEchoHandler()(*EchoHandler){
	return &EchoHandler{}
}
func (h *EchoHandler)Handle(ctx context.Context,conn net.Conn){
	if h.closing.Get() {
		_=conn.Close()
		return
	}

}