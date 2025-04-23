package config

import "github.com/alibaba/sentinel-golang/logging"

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
	SentinelInfo SentinelConfig `mapstructure:"sentinel" json:"sentinel"`
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

// SentinelConfig represent the general configuration of Sentinel.
type SentinelConfig struct {
	App struct {
		// Name represents the name of current running service.
		Name string `mapstructure:"name"`
		// Type indicates the classification of the service (e.g. web service, API gateway).
		Type int32 `mapstructure:"name"`
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
	// Logger indicates that using logger to replace default logging.
	Logger logging.Logger
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
