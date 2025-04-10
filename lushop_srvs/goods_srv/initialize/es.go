package initialize

import (
	"context"
	"fmt"
	"goodssrv/global"
	"goodssrv/model"
	"log"
	"os"

	"github.com/olivere/elastic/v7"
)

func Es() {
	// 初始化连接
	host := fmt.Sprintf("http://%s:%s", global.ServerConfig.EsInfo.Host, global.ServerConfig.EsInfo.Port)
	logger := log.New(os.Stdout, "lushop", log.LstdFlags)
	var err error
	global.EsClient, err = elastic.NewClient(elastic.SetURL(host), elastic.SetSniff(false),
		elastic.SetTraceLog(logger))
	if err != nil {
		panic(err)
	}

	// 新建mapping和index
	exist, err := global.EsClient.IndexExists(model.EsGoods{}.GetIndexName()).Do(context.Background())
	if err != nil {
		panic(err)
	}
	if !exist {
		global.EsClient.CreateIndex(model.EsGoods{}.GetIndexName()).BodyString(model.EsGoods{}.GetMapping()).Do(context.Background())
		if err != nil {
			panic(err)
		}
	}
}
