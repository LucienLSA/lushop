package config

type MySQLConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     string `mapstructure:"port" json:"port"`
	DbName   string `mapstructure:"db_name" json:"db_name"`
	User     string `mapstructure:"user" json:"user"`
	PassWord string `mapstructure:"password" json:"password"`
}

type ServerConfig struct {
	Name       string       `mapstructure:"name" json:"name"`
	MySQLInfo  MySQLConfig  `mapstructure:"mysql" json:"mysql"`
	ConsulInfo ConsulConfig `mapstructure:"consul" json:"consul"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port string `mapstructure:"port" json:"port"`
}
