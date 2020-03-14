package wechat

import "gopkg.in/chanxuehong/wechat.v2/mp/core"

var mux *core.ServeMux

func init()  {
	mux=core.NewServeMux()
	mux.DefaultEventHandleFunc(messageHandler)
	mux.DefaultMsgHandleFunc(eventHandler)
}



//消息处理器
func messageHandler(ctx *core.Context){

	Daemon.Logger.Info(ctx.MixedMsg)
}

//事件处理器
func eventHandler(ctx *core.Context){

	Daemon.Logger.Info(ctx.MixedMsg)

}

