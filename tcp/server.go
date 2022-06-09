package tcp

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)
type Config struct {
	Address string
	MaxConnect uint32
	Timeout time.Duration
}



func ListenAndServe(listener net.Listener,handler Handler,closeChan<-chan struct{}){
	//正常关闭
	go func() {
		<-closeChan
		log.Println("shutting down...")
		_=listener.Close()//停止监听，listener.Accept()会立即返回io.EOF
		_=handler.Close()//关闭应用层服务器
	}()
	//异常关闭,释放资源
	defer func(){
		_=listener.Close()
		_=handler.Close()
	}()
	ctx:=context.Background()
	var waitDone sync.WaitGroup
	for {
		//监听端口，阻塞直到收到新链接或者出现错误
		conn,err := listener.Accept()
		if err !=nil{
			break
		}
		//开启goroutine来处理新链接
		log.Println("accept link")
		waitDone.Add(1)
		go func(){
			defer func(){
				waitDone.Done()
			}()
			handler.Handle(ctx,conn)
		}()
		waitDone.Wait()
	}

}
func ListenAndServeWithSignal(cfg *Config,handler Handler)error{
	closeChan :=make(chan struct{})
	sigCh:=make(chan os.Signal)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig:=<-sigCh
		switch sig{
		case  syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			closeChan<- struct{}{}
		}
	}()
	listener,err :=net.Listen("tcp",cfg.Address)
	if err != nil{
		return err
	}
	ListenAndServe(listener,handler,closeChan)
	return nil
}