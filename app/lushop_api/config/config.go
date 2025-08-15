package config

type UserSrvConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type UserOpSrvConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type GoodsSrvConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type OrderSrvConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type InventorySrvConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port int    `mapstructure:"port" json:"port"`
	Name string `mapstructure:"name" json:"name"`
}

type OssConfig struct {
	ApiKey      string `mapstructure:"api_key" json:"api_key"`
	ApiSecrect  string `mapstructure:"api_secrect" json:"api_secrect"`
	Endpoint    string `mapstructure:"endpoint" json:"endpoint"`
	Bucket      string `mapstructure:"bucket" json:"bucket"`
	Host        string `mapstructure:"host" json:"host"`
	CallBackUrl string `mapstructure:"callback_url" json:"callback_url"`
	UploadDir   string `mapstructure:"upload_dir" json:"upload_dir"`
	ExpireTime  int64  `mapstructure:"expired_time" json:"expired_time"`
}
type JaegerConfig struct {
	ServiceName       string `mapstructure:"service_name" json:"service_name"`
	JaegerGinEndpoint string `mapstructure:"jaeger_gin_endpoint" json:"jaeger_gin_endpoint"`
}

type JwtConfig struct {
	Secret            string `mapstructure:"secret" json:"secret"`
	Key               string `mapstructure:"key" json:"key"`
	AccessExpireTime  int64  `mapstructure:"access_expired_time" json:"access_expired_time"`
	RefreshExpireTime int64  `mapstructure:"refresh_expired_time" json:"refresh_expired_time"`
}

type AliSmsConfig struct {
	ApiKey       string `mapstructure:"api_key" json:"api_key"`
	ApiSecret    string `mapstructure:"api_secret" json:"api_secret"`
	SignName     string `mapstructure:"sign_name" json:"sign_name"`
	PhoneNumber  string `mapstructure:"phone_number" json:"phone_number"`
	TemplateCode string `mapstructure:"template_code" json:"template_code"`
	RegionId     string `mapstructure:"region_id" json:"region_id"`
	Expire       int    `mapstructure:"expire" json:"expire"`
	CoolDown     int    `mapstructure:"cooldown" json:"cooldown"`
}

type ConsulConfig struct {
	Host string `mapstructure:"host" json:"host"`
	Port string `mapstructure:"port" json:"port"`
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

// SentinelConfig represent the general configuration of Sentinel.
type SentinelConfig struct {
	App struct {
		// Name represents the name of current running service.
		Name string `mapstructure:"name"`
		// Type indicates the classification of the service (e.g. web service, API gateway).
		Type int32 `mapstructure:"type"`
	}
	// Log represents configuration items related to logging.
	Log SentinelLogConfig `mapstructure:"log"`
	// Stat represents configuration items related to statistics.
	Stat StatConfig `mapstructure:"stat"`
	// UseCacheTime indicates whether to cache time(ms)
	UseCacheTime bool `mapstructure:"useCacheTime"`
}

// LogConfig represent the configuration of logging in Sentinel.
type SentinelLogConfig struct {
	// Dir represents the log directory path.
	Dir string `mapstructure:"dir"`
	// UsePid indicates whether the filename ends with the process ID (PID).
	UsePid bool `mapstructure:"usePid"`
	// Metric represents the configuration items of the metric log.
	Metric MetricLogConfig `mapstructure:"metric"`
}

// MetricLogConfig represents the configuration items of the metric log.
type MetricLogConfig struct {
	SingleFileMaxSize uint64 `mapstructure:"singleFileMaxSize"`
	MaxFileCount      uint32 `mapstructure:"maxFileCount"`
	FlushIntervalSec  uint32 `mapstructure:"flushIntervalSec"`
}

// StatConfig represents the configuration items of statistics.
type StatConfig struct {
	// GlobalStatisticSampleCountTotal and GlobalStatisticIntervalMsTotal is the per resource's global default statistic sliding window config
	GlobalStatisticSampleCountTotal uint32 `mapstructure:"globalStatisticSampleCountTotal"`
	GlobalStatisticIntervalMsTotal  uint32 `mapstructure:"globalStatisticIntervalMsTotal"`

	// MetricStatisticSampleCount and MetricStatisticIntervalMs is the per resource's default readonly metric statistic
	// This default readonly metric statistic must be reusable based on global statistic.
	MetricStatisticSampleCount uint32 `mapstructure:"metricStatisticSampleCount"`
	MetricStatisticIntervalMs  uint32 `mapstructure:"metricStatisticIntervalMs"`

	System SystemStatConfig `mapstructure:"system"`
}

// SystemStatConfig represents the configuration items of system statistics.
type SystemStatConfig struct {
	// CollectIntervalMs represents the collecting interval of the system metrics collector.
	CollectIntervalMs uint32 `mapstructure:"collectIntervalMs"`
}

type LogConfig struct {
	Level      string `mapstructure:"level" json:"level"`
	Filepath   string `mapstructure:"filepath" json:"filepath"`
	Filename   string `mapstructure:"filename" json:"filename"`
	MaxSize    int    `mapstructure:"max_size" json:"max_size"`
	MaxAge     int    `mapstructure:"max_age" json:"max_age"`
	MaxBackups int    `mapstructure:"max_backups" json:"max_backups"`
}

type ServerConfig struct {
	Name             string             `mapstructure:"name" json:"name"`
	Host             string             `mapstructure:"host" json:"host"`
	Port             int                `mapstructure:"port" json:"port"`
	Version          string             `mapstructure:"version" json:"version"`
	Tags             []string           `mapstructure:"tags" json:"tags"`
	JwtInfo          JwtConfig          `mapstructure:"jwt" json:"jwt"`
	AliSmsInfo       AliSmsConfig       `mapstructure:"ali_sms" json:"ali_sms"`
	RedisInfo        RedisConfig        `mapstructure:"redis" json:"redis"`
	SessionInfo      SessionConfig      `mapstructure:"session" json:"session"`
	ConsulInfo       ConsulConfig       `mapstructure:"consul" json:"consul"`
	JaegerInfo       JaegerConfig       `mapstructure:"jaeger" json:"jaeger"`
	OssInfo          OssConfig          `mapstructure:"oss" json:"oss"`
	AliPayInfo       AlipayConfig       `mapstructure:"alipay" json:"alipay"`
	UserSrvInfo      UserSrvConfig      `mapstructure:"user_srv" json:"user_srv"`
	UserOpSrvInfo    UserOpSrvConfig    `mapstructure:"userop_srv" json:"userop_srv"`
	GoodsSrvInfo     GoodsSrvConfig     `mapstructure:"goods_srv" json:"goods_srv"`
	OrderSrvInfo     OrderSrvConfig     `mapstructure:"order_srv" json:"order_srv"`
	InventorySrvInfo InventorySrvConfig `mapstructure:"inventory_srv" json:"inventory_srv"`
	SentinelInfo     SentinelConfig     `mapstructure:"sentinel" json:"sentinel"`
	LogInfo          LogConfig          `mapstructure:"log" json:"log"`
	NacosInfo        NacosConfig        `mapstructure:"nacos" json:"nacos"`
	Oauth2Info       Oauth2Config       `mapstructure:"oauth2" json:"oauth2"`
}

type SessionConfig struct {
	Name      string `mapstructure:"name" json:"name"`
	MaxAge    int    `mapstructure:"max_age" json:"max_age"`
	SecretKey string `mapstructure:"secret_key" json:"secret_key"`
}

type NacosConfig struct {
	Host      string `mapstructure:"host" json:"host"`
	Port      uint64 `mapstructure:"port" json:"port"`
	Namespace string `mapstructure:"namespace" json:"namespace"`
	User      string `mapstructure:"user" json:"user"`
	Password  string `mapstructure:"password" json:"password"`
	DataId    string `mapstructure:"dataid" json:"dataid"`
	Group     string `mapstructure:"group" json:"group"`
}

// Oauth2ScopeConfig 表示单个权限范围配置
// 对应yaml: oauth2.client.scope
// 例如：- id: all, title: "用户账号、手机、权限、角色等信息"
type Oauth2ScopeConfig struct {
	Id    string `mapstructure:"id" json:"id"`
	Title string `mapstructure:"title" json:"title"`
}

// Oauth2ClientConfig 表示单个oauth2客户端配置
// 对应yaml: oauth2.client
// 例如：- id: test_client_1, secret: test_secret_1, ...
type Oauth2ClientConfig struct {
	Id     string              `mapstructure:"id" json:"id"`
	Secret string              `mapstructure:"secret" json:"secret"`
	Name   string              `mapstructure:"name" json:"name"`
	Domain string              `mapstructure:"domain" json:"domain"`
	Scope  []Oauth2ScopeConfig `mapstructure:"scope" json:"scope"`
}

// Oauth2Config 表示oauth2整体配置
// 对应yaml: oauth2
// 包含access_token/refresh_token过期时间、签名key、客户端列表等
type Oauth2Config struct {
	AccessTokenExp  int                  `mapstructure:"access_token_exp" json:"access_token_exp"`   // access_token 过期时间（小时）
	RefreshTokenExp int                  `mapstructure:"refresh_token_exp" json:"refresh_token_exp"` // refresh_token 过期时间（小时）
	JwtSignedKey    string               `mapstructure:"jwt_signed_key" json:"jwt_signed_key"`       // JWT签名key
	Client          []Oauth2ClientConfig `mapstructure:"client" json:"client"`                       // 客户端列表
}
