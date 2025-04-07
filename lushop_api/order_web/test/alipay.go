package main

import (
	"fmt"

	"github.com/smartwalle/alipay/v3"
)

func main() {
	appID := "2021000147678720" // 后续改为硬编码
	privateKey := ""
	aliPublicKey := ""
	var client, err = alipay.New(appID, privateKey, false)
	err = client.LoadAliPayPublicKey(aliPublicKey)
	if err != nil {
		panic(err)
	}
	var p = alipay.TradePagePay{}
	p.NotifyURL = "http://127.0.0.1:8040/notify"
	p.ReturnURL = "http://127.0.0.1:8040/return"
	p.Subject = "lushop订单支付"
	p.OutTradeNo = "lucien_phone"
	p.TotalAmount = "10.00"
	p.ProductCode = "FAST_INSTANT_TRADE_PAY"

	url, err := client.TradePagePay(p)
	if err != nil {
		panic(err)
	}
	fmt.Println(url.String())
}
