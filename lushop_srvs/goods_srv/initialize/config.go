package initialize

import (
	"fmt"
	"goodssrv/global"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func Config() {
	// 从配置文件中读取配置
	mode := global.GetEnvInfoBool(global.Mode)
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("%s-pro.yaml", configFilePrefix)
	fmt.Println(mode)
	if mode {
		configFileName = fmt.Sprintf("%s-debug.yaml", configFilePrefix)
	}
	v := viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	// serverConfig := config.ServerConfig{}
	if err := v.Unmarshal(&global.ServerConfig); err != nil {
		panic(err)
	}
	zap.S().Info("配置信息:%v", global.ServerConfig)
	// fmt.Printf("配置信息:%v", global.ServerConfig)
	// 对象如何在其他文件中使用，因为它是一个局部变量，初始化全局变量
	// if err := v.Unmarshal(&global.NacosConfig); err != nil {
	// 	panic(err)
	// }
	// zap.S().Info("配置信息:%v", global.NacosConfig)
	fmt.Printf("%v", v.Get("name"))
	// 动态监控变化
	v.WatchConfig()
	v.OnConfigChange(func(in fsnotify.Event) {
		zap.S().Info("配置文件产生变化:%s", in.Name)
		_ = v.ReadInConfig()
		_ = v.Unmarshal(global.ServerConfig)
		zap.S().Info("配置信息:%v", global.ServerConfig)
	})

	// // 从nacos中读取配置信息
	// sc := []constant.ServerConfig{
	// 	{
	// 		IpAddr: global.NacosConfig.NacosInfo.Host,
	// 		Port:   global.NacosConfig.NacosInfo.Port,
	// 	},
	// }
	// cc := constant.ClientConfig{
	// 	NamespaceId:         global.NacosConfig.NacosInfo.Namespace, // 如果需要支持多namespace，我们可以创建多个client,它们有不同的NamespaceId。当namespace是public时，此处填空字符串。
	// 	TimeoutMs:           5000,
	// 	NotLoadCacheAtStart: true,
	// 	LogDir:              "./temp/nacos/log",
	// 	CacheDir:            "./temp/nacos/cache",
	// 	LogLevel:            "debug",
	// }
	// clientConfig, err := clients.CreateConfigClient(map[string]interface{}{
	// 	"serverConfigs": sc,
	// 	"clientConfig":  cc,
	// })
	// if err != nil {
	// 	panic(err)
	// }
	// content, err := clientConfig.GetConfig(vo.ConfigParam{
	// 	DataId: global.NacosConfig.NacosInfo.DataId,
	// 	Group:  global.NacosConfig.NacosInfo.Group})
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(global.NacosConfig.NacosInfo.DataId, global.NacosConfig.NacosInfo.Group)
	// fmt.Println(content)
	// // serverConfig := config.ServerConfig{} // 局部变量
	// err = json.Unmarshal([]byte(content), &global.ServerConfig)
	// if err != nil {
	// 	zap.S().Fatalf("读取nacos配置失败:%s", err.Error())
	// }
	// fmt.Println(&global.ServerConfig)
	// err = clientConfig.ListenConfig(vo.ConfigParam{
	// 	DataId: global.NacosConfig.NacosInfo.DataId,
	// 	Group:  global.NacosConfig.NacosInfo.Group,
	// 	OnChange: func(namespace, group, dataId, data string) {
	// 		fmt.Println("配置文件产生变化")
	// 		fmt.Println("group:" + group + ", dataId:" + dataId + ", data:" + data)
	// 	},
	// })
	// time.Sleep(1 * time.Second)
}
