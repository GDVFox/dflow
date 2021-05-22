package config

import (
	"github.com/GDVFox/dflow/machine_node/watcher"
	"github.com/GDVFox/dflow/util"
	"github.com/GDVFox/dflow/util/httplib"
	"github.com/GDVFox/dflow/util/storage"
)

// Conf глобальный конфиг синглтон.
var Conf = NewConfig()

// Config конфигурация сервиса.
type Config struct {
	HTTP    *httplib.HTTPConfig `yaml:"http"`
	Logging *util.LoggingConfig `yaml:"logging"`
	ETCD    *storage.ETCDConfig `yaml:"etcd"`
	Watcher *watcher.Config     `yaml:"watcher"`
	// ActionStartRetry настройки ретрая для признания действия
	// запущенным. Если рантайм не ответил за N попыток, то считается
	// что он не запущен.
	ActionStartRetry *util.RetryConfig `yaml:"action-start-retry"`
	// RuntimePath путь к бинарному файлу рантайма.
	RuntimePath string `yaml:"runtime-path"`
	// RuntimeLogsDir путь к директории, в которой будут хранится логи
	// запущенных рантаймов.
	RuntimeLogsDir string `yaml:"runtime-logs-dir"`
	// RuntimeLogsLevel уровень логирования в загруженных рантаймах
	RuntimeLogsLevel string `yaml:"runtime-logs-level"`
}

// NewConfig создает конфиг с настройками по-умолчанию
func NewConfig() *Config {
	return &Config{
		HTTP:             httplib.NewtHTTPConfig(),
		Logging:          util.NewLoggingConfig(),
		ETCD:             storage.NewETCDConfig(),
		Watcher:          watcher.NewConfig(),
		ActionStartRetry: util.NewRetryConfig(),
		RuntimePath:      "runtime",
		RuntimeLogsDir:   "runtime-logs",
		RuntimeLogsLevel: "info",
	}
}
