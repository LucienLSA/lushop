package config

type UserSrvConfig struct {
	Host    string `mapstructure:"host" json:"host"`
	Port    int    `mapstructure:"port" json:"port"`
	Version string `mapstructure:"version" json:"version"`
	Name    string `mapstructure:"name" json:"name"`
}

type JwtConfig struct {
	SigningKey string `mapstructure:"key" json:"key"`
	ExpireTime int64  `mapstructure:"expired_time" json:"expired_time"`
}

type AliSmsConfig struct {
	ApiKey       string `mapstructure:"api_key" json:"api_key"`
	ApiSecrect   string `mapstructure:"api_secret" json:"api_secret"`
	SignName     string `mapstructure:"sign_name" json:"sign_name"`
	PhoneNumber  string `mapstructure:"phone_number" json:"phone_number"`
	TemplateCode string `mapstructure:"template_code" json:"template_code"`
	RegionId     string `mapstructure:"region_id" json:"region_id"`
	Expire       int    `mapstructure:"expire" json:"expire"`
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

type ServerConfig struct {
	Name        string        `mapstructure:"name" json:"name"`
	Port        int           `mapstructure:"port" json:"port"`
	Mode        string        `mapstructure:"mode" json:"mode"`
	UserSrvInfo UserSrvConfig `mapstructure:"user_srv" json:"user_srv"`
	JwtInfo     JwtConfig     `mapstructure:"jwt" json:"jwt"`
	AliSmsInfo  AliSmsConfig  `mapstructure:"ali_sms" json:"ali_sms"`
	RedisInfo   RedisConfig   `mapstructure:"redis" json:"redis"`
	ConsulInfo  ConsulConfig  `mapstructure:"consul" json:"consul"`
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
