package main

import (
	"fmt"
	"lushopapi/global"

	"github.com/smartwalle/alipay/v3"
)

func main() {
	appID := global.GetEnvInfoStr(global.ServerConfig.AliPayInfo.AppID)
	privateKey := global.GetEnvInfoStr(global.ServerConfig.AliPayInfo.PrivateKey)
	aliPublicKey := global.GetEnvInfoStr(global.ServerConfig.AliPayInfo.AliPublicKey)
	aliPublicKey := "MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAkr5LOUbUzaJ+O1kzRyTLhWdQFNc3xzwRkCxUQr0KXgH+ah0l2y7lc+TJBI/mv3G8wtdNarqWHmb9v1Um3Da8Uila9WEGaa2QlvcZrwkvrkl8lNjbNbjyoG2Rak3AstrQnZnuzxBC2s6L0a/TFNMvSDg/AP1igcqmDvswQCcQkUdQiILofL/AAMJ+6VA/koXEBRTZkjboOI5wBXtupneIQ2KsI2AE3EagfxPKJdAMNrkobKyknI/2IbFeXV4M7SxzWJo8IZlPSeiu+wTGq0aurfJb3IAHQzcM6PEnALPeepc7VwvgygjUkmnwWk8GJyhOwkcoN3eigyAltXQnLCTuvQIDAQAB"
	var client, err = alipay.New(appID, privateKey, false)
	err = client.LoadAliPayPublicKey(aliPublicKey)
	if err != nil {
		panic(err)
	}
	var p = alipay.TradePagePay{}
	p.NotifyURL = "http://127.0.0.1:8101/alipay/notify"
	// p.ReturnURL = "http://127.0.0.1:8000/return"
	p.Subject = "lushop订单支付"
	p.OutTradeNo = "lucien_computer"
	p.TotalAmount = "11.00"
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	url, err := client.TradePagePay(p)
	if err != nil {
		panic(err)
	}
	fmt.Println(url.String())
}
