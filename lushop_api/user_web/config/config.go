package config

type UserSrvConfig struct {
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
	Mode    string `mapstructure:"mode"`
	Version string `mapstructure:"version"`
}

type JwtConfig struct {
	SigningKey string `mapstructure:"key"`
	ExpireTime int64  `mapstructure:"expired_time"`
}

type AliSmsConfig struct {
	ApiKey       string `mapstructure:"api_key"`
	ApiSecrect   string `mapstructure:"api_secret"`
	SignName     string `mapstructure:"sign_name"`
	PhoneNumber  string `mapstructure:"phone_number"`
	TemplateCode string `mapstructure:"template_code"`
	RegionId     string `mapstructure:"region_id"`
	Expire       int    `mapstructure:"expire"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Password string `mapstructure:"password"`
	PoolSize int    `mapstructure:"pool_size"`
	DB       int    `mapstructure:"db"`
}

type ServerConfig struct {
	Name        string        `mapstructure:"name"`
	Port        int           `mapstructure:"port"`
	UserSrvInfo UserSrvConfig `mapstructure:"user_srv"`
	JwtInfo     JwtConfig     `mapstructure:"jwt"`
	AliSmsInfo  AliSmsConfig  `mapstructure:"ali_sms"`
	RedisInfo   RedisConfig   `mapstructure:"redis"`
}
