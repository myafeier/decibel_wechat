package wechat

import (
	"io/ioutil"
	"testing"
)

var weappConfig = WeChatMicroAppConfig{
	AppId:  "wx884199c1f98151f3",
	Secret: "8119bb3228e1b80705f94971551301f7",
}

func init() {
	NewMicroappDaemon(&weappConfig)
}

func TestMicroappServiceGenerateWxappCode(t *testing.T) {

	data, err := MicroappDaemon.GenerateWxappCode("pages/index/index", "code=xx", 400)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = ioutil.WriteFile("./wxappcode.png", data, 0744)
	if err != nil {
		t.Fatal(err.Error())
	}

}
