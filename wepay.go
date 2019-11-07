// 微信支付
package decibel_wechat

import (
	"crypto/md5"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"gopkg.in/chanxuehong/wechat.v2/mch/core"
	"gopkg.in/chanxuehong/wechat.v2/mch/pay"
	"github.com/go-xorm/xorm"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"
)

type WePay struct {
	AppId  string
	config *WePayConfig
	db     *xorm.Session
	logger ILogger
	mutex  sync.Mutex
}
// 统一下单，返回再次签名的json数据
func (self *WePay) UnifiedOrder(userOpenId,orderSn, orderInfo string,localIp string, amount int) (jsonByte []byte, err error) {
	tradeType := "JSAPI"
	nonceStr := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%s%d", orderSn, time.Now().Unix()))))
	reqMap := make(map[string]string)
	reqMap["openid"] = userOpenId
	reqMap["nonce_str"] = nonceStr
	reqMap["body"] = orderInfo
	reqMap["out_trade_no"] = orderSn
	reqMap["total_fee"] = fmt.Sprintf("%d", amount)
	reqMap["spbill_create_ip"] = localIp
	reqMap["notify_url"] = self.config.NotifyUrl
	reqMap["trade_type"] = tradeType
	reqMap["sign_type"] = "MD5"

	self.logger.Debug(reqMap)

	client := core.NewClient(self.AppId, self.config.MerchantId, self.config.MerchantSecret, nil)
	response, err := pay.UnifiedOrder(client, reqMap)
	if err != nil {
		self.logger.Error(err)
		return
	}

	if response["return_code"] != core.ResultCodeSuccess {
		err = fmt.Errorf("unified order error,code:%s,des:%s", response["err_code"], response["err_code_des"])
		self.logger.Error(err)
		return
	}

	payInfo:=&WePayment{}
	payInfo.Stat=WePayStatToPay
	payInfo.NotifyUrl=reqMap["notify_url"]
	payInfo.Openid=reqMap["openid"]
	payInfo.NonceStr=reqMap["nonce_str"]
	payInfo.Body=reqMap["body"]
	payInfo.OutTradeNo=reqMap["out_trade_no"]
	payInfo.TotalFee=amount
	payInfo.SpbillCreateIp=reqMap["spbill_create_ip"]
	payInfo.TradeType=reqMap["trade_type"]
	payInfo.SignType=reqMap["sign_type"]
	payInfo.PrepayId= response["prepay_id"]
	payInfo.TransactionId= response["TransactionId"]
	_,err=self.db.Insert(payInfo)
	if err != nil {
		self.logger.Error(err)
		return
	}
	return self.signAgain(payInfo.NonceStr,payInfo.PrepayId,amount)
}

func (self *WePay) signAgain(nonceStr, prepayId string, amount int) (signInfo []byte, err error) {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	packageStr := fmt.Sprintf("prepay_id=%s", prepayId)
	signType := "MD5"
	signAgainData := map[string]string{
		"appId":     self.AppId,
		"timeStamp": timestamp,
		"nonceStr":  nonceStr,
		"package":   packageStr,
		"signType":  signType,
	}

	str := fmt.Sprintf("appId=%s&nonceStr=%s&package=%s&signType=%s&timeStamp=%s&key=%s", self.AppId, nonceStr, packageStr, signType, timestamp, self.config.MerchantSecret)
	signAgainData["paySign"] = strings.ToUpper(fmt.Sprintf("%x", md5.Sum([]byte(str))))
	if err != nil {
		self.logger.Error(err)
		return
	}
	signInfo, err = json.Marshal(signAgainData)
	return
}

// 支付成功后的回调
//  watcher ，支付成功观察者
func (self *WePay) CallBack(request io.Reader,watcher IPayCallbackWatcher) (err error) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	defer func() {
		if err != nil {
			self.logger.Error(err)
			return
		}
	}()
	data, err := decodeXMLToMap(request)

	hadSign, ok := data["sign"]
	if !ok {
		return
	}
	wantSign := core.Sign(data, self.config.MerchantSecret, nil)
	if hadSign!=wantSign {
		err = fmt.Errorf("sign mismatch,\nrequest sign: %s\n want sign: %s", hadSign, wantSign)
		return
	}
	var orderSn,transactionId string
	var orderAmount int
	var payStatus bool

	switch data["result_code"] {
	case "SUCCESS":
		orderSn, ok = data["out_trade_no"]
		if !ok {
			err=fmt.Errorf("request without out_trade_no")
			return
		}
		transactionId, ok = data["transaction_id"]
		if !ok {
			err=fmt.Errorf("request without transaction_id")
			return
		}
		totalFee, ok := data["total_fee"]
		if !ok {
			err=fmt.Errorf("request without total_fee")
			return
		}
		orderAmount, err = strconv.Atoi(totalFee)
		if err != nil {
			return
		}
		payStatus=true

	case "FAIL":
		orderSn, ok = data["out_trade_no"]
		if !ok {
			err=fmt.Errorf("request without out_trade_no")
			return
		}
		payStatus=false

	default:
		err=fmt.Errorf("WEPAY Notify Error: Request Body:%+v\n", data)
		return
	}

	wxOrder:=new(WePayment)
	var has bool
	has,err=self.db.Where("out_trade_no=?",orderSn).Get(wxOrder)
	if err != nil {
		return
	}
	if !has{
		err=fmt.Errorf("call back order not exist,out_trade_no:%s",orderSn)
		return
	}
	if wxOrder.TotalFee!=orderAmount{
		err=fmt.Errorf("call back order amount not eqal,orderSn:%s,callback:%d, localAmount: %d",orderSn,orderAmount,wxOrder.TotalFee)
		return
	}
	if wxOrder.Stat!=WePayStatToPay{
		err=fmt.Errorf("call back order already deal! orderSN: %s ",orderSn)
		return
	}

	if payStatus{

		wxOrder.Stat=WePayStatPaySuccess
		wxOrder.TransactionId+=","+transactionId
		_,err=self.db.ID(wxOrder.Id).Cols("stat,return_code,transaction_id").Update(wxOrder)
		if err != nil {
			return
		}
		err=watcher.OrderPaySuccess(orderSn,orderAmount)
		if err != nil {
			self.logger.Error(err)
			wxOrder.WatcherResult=err.Error()
		}else{
			wxOrder.WatcherResult="success"
		}


	}else{
		wxOrder.Stat=WePayStatPayFail
		wxOrder.ReturnCode=data["err_code"]+"|"+data["err_code_des"]
		self.logger.Info(data)

		_,err=self.db.ID(wxOrder.Id).Cols("stat,return_code").Update(wxOrder)

		err=watcher.OrderPayFail(orderSn,orderAmount)
		if err != nil {
			self.logger.Error(err)
			wxOrder.WatcherResult=err.Error()
		}else{
			wxOrder.WatcherResult="success"
		}

	}

	_,err=self.db.ID(wxOrder.Id).Cols("watcher_result").Update(wxOrder)
	if err != nil {
		return
	}
	return
}

func decodeXMLToMap(r io.Reader) (m map[string]string, err error) {
	m = make(map[string]string)
	var (
		decoder = xml.NewDecoder(r)
		depth   = 0
		token   xml.Token
		key     string
		value   strings.Builder
	)
	for {
		token, err = decoder.Token()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}

		switch v := token.(type) {
		case xml.StartElement:
			depth++
			switch depth {
			case 2:
				key = v.Name.Local
				value.Reset()
			case 3:
				if err = decoder.Skip(); err != nil {
					return
				}
				depth--
				key = "" // key == "" indicates that the node with depth==2 has children
			}
		case xml.CharData:
			if depth == 2 && key != "" {
				value.Write(v)
			}
		case xml.EndElement:
			if depth == 2 && key != "" {
				m[key] = value.String()
			}
			depth--
		}
	}
}
