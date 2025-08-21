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
