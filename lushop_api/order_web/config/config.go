package config

type OrderSrvConfig struct {
	Host    string `mapstructure:"host" json:"host"`
	Port    int    `mapstructure:"port" json:"port"`
	Version string `mapstructure:"version" json:"version"`
	Name    string `mapstructure:"name" json:"name"`
}

type GoodsSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}
type InventorySrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type JwtConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
	ExpireTime int64  `mapstructure:"expired_time" json:"expired_time"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     string `mapstructure:"port" json:"port"`
	Password string `mapstructure:"password" json:"password"`
	PoolSize int    `mapstructure:"pool_size" json:"pool_size"`
	DB       int    `mapstructure:"db" json:"db"`
}

type AlipayConfig struct {
	AppID        string `mapstructure:"app_id" json:"app_id"`
	PrivateKey   string `mapstructure:"private_key" json:"private_key"`
	AliPublicKey string `mapstructure:"ali_public_key" json:"ali_public_key"`
	NotifyURL    string `mapstructure:"notify_url" json:"notify_url"`
	ReturnURL    string `mapstructure:"return_url" json:"return_url"`
	ProductCode  string `mapstructure:"product_code" json:"product_code"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port string `mapstructure:"port" json:"port"`
}

type ServerConfig struct {
	Name             string             `mapstructure:"name" json:"name"`
	Host             string             `mapstructure:"host" json:"host"`
	Tags             []string           `mapstructure:"tags" json:"tags"`
	Port             int                `mapstructure:"port" json:"port"`
	OrderSrvInfo     OrderSrvConfig     `mapstructure:"order_srv" json:"order_srv"`
	GoodsSrvInfo     GoodsSrvConfig     `mapstructure:"goods_srv" json:"goods_srv"`
	InventorySrvInfo InventorySrvConfig `mapstructure:"inventory_srv" json:"inventory_srv"`
	JwtInfo          JwtConfig          `mapstructure:"jwt" json:"jwt"`
	ConsulInfo       ConsulConfig       `mapstructure:"consul" json:"consul"`
	AlipayInfo       AlipayConfig       `mapstructure:"alipay" json:"alipay"`
}

type NacosConfig struct {
	NacosInfo NacosInfo `mapstructure:"nacos"`
}

type NacosInfo struct {
	Host      string `mapstructure:"host"`
	Port      uint64 `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	User      string `mapstructure:"user"`
	Password  string `mapstructure:"password"`
	DataId    string `mapstructure:"dataid"`
	Group     string `mapstructure:"group"`
}
