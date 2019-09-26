package decibel_wechat

type Config struct {
	*WeChatMicroAppConfig
	*WePayConfig
}

type WeChatMicroAppConfig struct {
	OriginId     string //小程序原始id
	AppId        string //
	Secret       string
	Base64AESKey string
}

type WePayConfig struct {
	MerchantId     string //商户id
	MerchantSecret string //商户
	NotifyUrl      string //回调地址
}
