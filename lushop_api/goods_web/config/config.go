package config

type GoodsSrvConfig struct {
	Host    string `mapstructure:"host" json:"host"`
	Port    int    `mapstructure:"port" json:"port"`
	Version string `mapstructure:"version" json:"version"`
	Name    string `mapstructure:"name" json:"name"`
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

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port string `mapstructure:"port" json:"port"`
}

type JaegerConfig struct {
	ServiceName       string `mapstructure:"service_name" json:"service_name"`
	JaegerGinEndpoint string `mapstructure:"jaeger_gin_endpoint" json:"jaeger_gin_endpoint"`
}

type LogConfig struct {
	Level      string `mapstructure:"level" json:"level"`
	FilePath   string `mapstructure:"filepath" json:"filepath"`
	FileName   string `mapstructure:"filename" json:"filename"`
	MaxSize    int    `mapstructure:"max_size" json:"max_size"`
	MaxAge     int    `mapstructure:"max_age" json:"max_age"`
	MaxBackUps int    `mapstructure:"max_backups" json:"max_backups"`
}

type ServerConfig struct {
	Name         string         `mapstructure:"name" json:"name"`
	Host         string         `mapstructure:"host" json:"host"`
	Tags         []string       `mapstructure:"tags" json:"tags"`
	Port         int            `mapstructure:"port" json:"port"`
	GoodsSrvInfo GoodsSrvConfig `mapstructure:"goods_srv" json:"goods_srv"`
	JwtInfo      JwtConfig      `mapstructure:"jwt" json:"jwt"`
	ConsulInfo   ConsulConfig   `mapstructure:"consul" json:"consul"`
	JaegerInfo   JaegerConfig   `mapstructure:"jaeger" json:"jaeger"`
	LogInfo      LogConfig      `mapstructure:"log" json:"log"`
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
