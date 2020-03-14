package wechat

type ILogger interface {
	Debug(...interface{})
	Error(...interface{})
	Info(...interface{})
}

// 订单支付
type IPayCallbackWatcher interface {
	//支付成功
	OrderPaySuccess(orderSn string, Amount int) error

	//支付失败
	OrderPayFail(orderSn string, Amount int) error
}
