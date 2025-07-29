package initialize

import (
	"context"
	"fmt"
	"goodssrv/global"
	"goodssrv/model"
	"log"
	"os"

	"github.com/olivere/elastic/v7"
	"go.uber.org/zap"
)

func Es() {
	// 初始化连接
	host := fmt.Sprintf("http://%s:%s", global.ServerConfig.EsInfo.Host, global.ServerConfig.EsInfo.Port)
	logger := log.New(os.Stdout, "lushop", log.LstdFlags)
	var err error
	global.EsClient, err = elastic.NewClient(elastic.SetURL(host), elastic.SetSniff(false),
		elastic.SetTraceLog(logger))
	// elasticsearch Go 客户端在执行每个 HTTP 请求时，把详细的 trace 日志（如请求内容、响应内容、耗时等）通过你自定义的 logger 输出出来
	if err != nil {
		panic(err)
	}

	//查询index是否存在
	exist, err := global.EsClient.IndexExists(model.EsGoods{}.GetIndexName()).Do(context.Background())
	if err != nil {
		panic(err)
	}
	// 不存在就新建mapping和index
	if !exist {
		_, err2 := global.EsClient.CreateIndex(model.EsGoods{}.GetIndexName()).BodyString(model.EsGoods{}.GetMapping()).Do(context.Background())
		if err2 != nil {
			zap.S().Fatalf("创建索引%s失败:%s", model.EsGoods{}.GetIndexName(), err2.Error())
		}
	}
}
