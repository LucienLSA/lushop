package config

type MySQLConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     string `mapstructure:"port" json:"port"`
	DbName   string `mapstructure:"db_name" json:"db_name"`
	User     string `mapstructure:"user" json:"user"`
	PassWord string `mapstructure:"password" json:"password"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     string `mapstructure:"port" json:"port"`
	Password string `mapstructure:"password" json:"password"`
	PoolSize int    `mapstructure:"pool_size" json:"pool_size"`
	DB       int    `mapstructure:"db" json:"db"`
}

type GoodsSrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type InventorySrvConfig struct {
	Name string `mapstructure:"name" json:"name"`
}

type ServerConfig struct {
	Name         string         `mapstructure:"name" json:"name"`
	Host         string         `mapstructure:"host" json:"host"`
	Tags         []string       `mapstructure:"tags" json:"tags"`
	Port         int            `mapstructure:"port" json:"port"`
	MySQLInfo    MySQLConfig    `mapstructure:"mysql" json:"mysql"`
	ConsulInfo   ConsulConfig   `mapstructure:"consul" json:"consul"`
	RedisInfo    RedisConfig    `mapstructure:"redis" json:"redis"`
	JaegerInfo   JaegerConfig   `mapstructure:"jaeger" json:"jaeger"`
	RocketMQInfo RocketMQConfig `mapstructure:"rocketmq" json:"rocketmq"`

	GoodsSrvInfo     GoodsSrvConfig     `mapstructure:"goods_srv" json:"goods_srv"`
	InventorySrvInfo InventorySrvConfig `mapstructure:"inventory_srv" json:"inventory_srv"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port string `mapstructure:"port" json:"port"`
}

type NacosConfig struct {
	NacosInfo NacosInfo `mapstructure:"nacos"`
}

type JaegerConfig struct {
	Host        string `mapstructure:"host" json:"host"`
	Port        string `mapstructure:"port" json:"port"`
	ServiceName string `mapstructure:"service_name" json:"service_name"`
	TracerName  string `mapstructure:"tracer_name" json:"tracer_name"`
}
type RocketMQConfig struct {
	Host               string `mapstructure:"host" json:"host"`
	Port               string `mapstructure:"port" json:"port"`
	TopicReback        string `mapstructure:"topic_reback" json:"topic_reback"`
	TopicTimeOut       string `mapstructure:"topic_timeout" json:"topic_timeout"`
	ConsumerGroup          string `mapstructure:"consumer_group" json:"consumer_group"`
	ProducerGroupOrder string `mapstructure:"producer_group_order" json:"producer_group_order"`
	ProducerGroupInventory  string `mapstructure:"producer_group_inventory" json:"producer_group_inventory"`
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
