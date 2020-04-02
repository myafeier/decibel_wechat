package wechat

type Config struct {
	*WeChatMicroAppConfig
	*WePayConfig
	*WePayVendorConfig
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
type WePayVendorConfig struct {
	AppId          string //服务商appId
	MerchantId     string //商户id
	MerchantSecret string //商户
}
