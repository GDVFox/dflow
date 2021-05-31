package config

import (
	"github.com/GDVFox/dflow/meta_node/watcher"
	"github.com/GDVFox/dflow/util"
	"github.com/GDVFox/dflow/util/httplib"
	"github.com/GDVFox/dflow/util/storage"
)

// Conf глобальный конфиг синглтон.
var Conf = NewConfig()

// Config конфигурация сервиса.
type Config struct {
	HTTP    *httplib.HTTPConfig        `yaml:"http"`
	Logging *util.LoggingConfig        `yaml:"logging"`
	ETCD    *storage.ETCDConfig        `yaml:"etcd"`
	Watcher *watcher.PlanWatcherConfig `yaml:"watcher"`
}

// NewConfig создает конфиг с настройками по-умолчанию
func NewConfig() *Config {
	return &Config{
		HTTP:    httplib.NewtHTTPConfig(),
		Logging: util.NewLoggingConfig(),
		ETCD:    storage.NewETCDConfig(),
		Watcher: watcher.NewPlanWatcherConfig(),
	}
}
