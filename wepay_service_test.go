package wechat

import "testing"

const (
	VendorCert = `-----BEGIN CERTIFICATE-----
MIID6TCCAtGgAwIBAgIUe4qOmW31RSEtigL7uTjiZU3e5GMwDQYJKoZIhvcNAQEL
BQAwXjELMAkGA1UEBhMCQ04xEzARBgNVBAoTClRlbnBheS5jb20xHTAbBgNVBAsT
FFRlbnBheS5jb20gQ0EgQ2VudGVyMRswGQYDVQQDExJUZW5wYXkuY29tIFJvb3Qg
Q0EwHhcNMjAwMjI3MDcwMDMwWhcNMjUwMjI1MDcwMDMwWjB7MRMwEQYDVQQDDAox
NTc3ODIwMTcxMRswGQYDVQQKDBLlvq7kv6HllYbmiLfns7vnu58xJzAlBgNVBAsM
HuaYhuaYjuWIhui0neenkeaKgOaciemZkOWFrOWPuDELMAkGA1UEBgwCQ04xETAP
BgNVBAcMCFNoZW5aaGVuMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA
27hqGlDjud9tvJM20JkC4dHqMfl9WlXpXhB2nWF9VQ4Ly/3svQ/VMnVb+eX0j+bf
H1b1dDiYdlkl2pCWfeZL7KY75+LiuRzD2M+DrcvZfVsXR6aZHVmq7W6TeJFyX0SF
PIpP+7Ta5WqhinqU3buUUoJe0kK/gLhSZN33mbZdBASh7eHqf98p5ImJFpoOhMop
zYaRhE3VFiPK+4cKT96jKp2gFXm6Op/kg6/WVTGTlDfP7y2CnFVJpCDjj5R81otM
yDTxknOc6PX/JevSPzZA5vcrWXz9jVUtswdOktIgylUAef2dcBaRjcT0hdhJLBzq
Lb+bRaSHsYmkNKX7WSajOwIDAQABo4GBMH8wCQYDVR0TBAIwADALBgNVHQ8EBAMC
BPAwZQYDVR0fBF4wXDBaoFigVoZUaHR0cDovL2V2Y2EuaXRydXMuY29tLmNuL3B1
YmxpYy9pdHJ1c2NybD9DQT0xQkQ0MjIwRTUwREJDMDRCMDZBRDM5NzU0OTg0NkMw
MUMzRThFQkQyMA0GCSqGSIb3DQEBCwUAA4IBAQCUvOWJNW2aXza/7cpTR+2ADvx8
WbyMRec3yRvo//UyGZI+ylQ+30pGBo51yqreQEu7sARuMmsQ3oRBFdhldvl0jAYl
awYZS1z4nRQDq4xlhak2akwguPLxdOVHq848c1z2uZxSf2OKx/INQDj/+q4iQ7WM
i+EX6UvSEZ8g6p0y7UHuJdAbDw/EvRtm+wr854EnjPS2DIIGoi/i0kNzS0DnZFwh
897Xb5mhfcHr4OJaHQFMvb5iHIidOTKJqeIxhyhpSlLPuXPFea7EM/B1X0sXrj2L
zvmbN3Zd3Sv2tFLcGa2upiCwkjNx2CJmWvlCdfcioSbP+I4ZiIfH+EUGP7Pu
-----END CERTIFICATE-----`
	VendorKey = `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDbuGoaUOO53228
kzbQmQLh0eox+X1aVeleEHadYX1VDgvL/ey9D9UydVv55fSP5t8fVvV0OJh2WSXa
kJZ95kvspjvn4uK5HMPYz4Oty9l9WxdHppkdWartbpN4kXJfRIU8ik/7tNrlaqGK
epTdu5RSgl7SQr+AuFJk3feZtl0EBKHt4ep/3ynkiYkWmg6EyinNhpGETdUWI8r7
hwpP3qMqnaAVebo6n+SDr9ZVMZOUN8/vLYKcVUmkIOOPlHzWi0zINPGSc5zo9f8l
69I/NkDm9ytZfP2NVS2zB06S0iDKVQB5/Z1wFpGNxPSF2EksHOotv5tFpIexiaQ0
pftZJqM7AgMBAAECggEAaUPy4WS7lYNrqZrF+i40aUgOcZ7b5XmfcodkrIXWa/ds
w3CGCRYayC/dvt3iy0aKwxMASYwLzzdqoUoAL5Uz7s69iJz0jkcvtSHGLm+pZRtN
DfDNDni1IUeGs47LQsUrKBQDuc2tyZfKiPOteoWxy83V69o6sUqdfuGxDB2IdrF2
G2PtpLS6lWUmMMMWZgU/LyCpnIRzC96zks9Ug4lJf1+ccU2L9XoET/vSoaIDVaRX
hE6brAgFzrCwyZjsLSDp5K/EFRlRx0fT2t7UNdeGKTMTRLM/R/OjGRZsw3+/kii/
fJx/ltKuwOmsuolprthYEIeECLrPPNNWMy4ndSss0QKBgQD1C/bJ3tie2jkCoryn
L03EtqNsi6XU/mYIiWLkcg7MfTuRlnbndsCO6jtt7EYLcOHD/66T/z5/IWUvC3Ba
R2b0pmfRhvnAbgT1Mjdy0O5ALV84hEXJEYutBcKKDUc5jUIg5S7RDvMX8fwFOlG6
V2Ytt4HhEUKOnswXulmKzFFBBQKBgQDliqTzdD5/WOpVhEZ9x4F5B49nMyXswBP+
ilEYSsqpMTsFDAxilR7/awtnTjAgYYALHWooTUDpqsdR38meMlj5ZJi5qfM5s38k
o/V6Fj5AWRzUIQm1kxDR5RtdsRgDatxMKDBDJa6PFn1nmsqrZfU/uWmBppak4mLb
wc0HlsaHPwKBgQDSrG6IL+bc65CIC5FVyv15Wew2rfjsnarrO/KhpM3EUQadrFad
uSLju81MPA4cV/hBodhdtNvuQK/VOmhltW12eHpZUUn3fp8Ujw/MzoOG+XscA9xb
eZI0NveB6NiLSj7IOUF+yvOEaq7Zb8JECk/2jgZDkas/Ipck7zl8cxyIYQKBgQCU
YPcX2MC9mUCBXywiCmELV3O/hjSxwcgq9kZNqasvi39XV955q2OKQCvy73v0spIO
nUkOHEIlyhtmNX8jH/Cb5gdDnTR4zCsYCFSaQt6iwff8uA6KrTJmO+9gtSWMr/sP
z7rC7QzVuff+jPUNrq7GLpihEoq2sxCsda6PhUt4CQKBgEjGOftdJ99lXzIVvMoq
NNuifI8/1uWwGhQu+4pgDKhd2WOt6aVxwR0fUcHwYaE493EypQvKcLsfZUsmVAQQ
1VlRPCC7QHoiBLHhcSu5FFnJAshShXEEmDtoGDSbOPfSvDTY796aPNo6A+B0/Yan
EL6vYKgn/A0S/pF+GRFPsl3a
-----END PRIVATE KEY-----`
)

func TestWepayServiceRedpackSend(t *testing.T) {
	serv := &WePayService{
		AppId: "wxd588716ce979c102",
		config: &WePayConfig{
			MerchantId:     "1582568201",
			MerchantSecret: "Uydh7635ysgh89ikojs63526352fhjdk",
			NotifyUrl:      "https://z1.u1200.com/api/v1/notify/wepay",
		},
		vendor: &WePayVendorConfig{
			MerchantId:     "1577947881",
			MerchantSecret: "Dushuu76537dyUydj987676jhyyu7ddd",
			AppId:          "wx7efbb7380e6f5fb2",
		},
		logger:   NewDefaultLogger(),
		CertData: []byte(VendorCert),
		KeyData:  []byte(VendorKey),
	}
	url, err := serv.RedpackSend(1, 1, "test", 1)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("url: %s", url)

}
