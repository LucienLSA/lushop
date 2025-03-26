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

type ServerConfig struct {
	Name        string        `mapstructure:"name"`
	Port        int           `mapstructure:"port"`
	UserSrvInfo UserSrvConfig `mapstructure:"user_srv"`
	JwtInfo     JwtConfig     `mapstructure:"jwt"`
}
