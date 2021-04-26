package logs

import "github.com/GDVFox/dflow/util"

// Logger синлтон-объект для логирования.
var Logger *util.Logger

// InitLogger инициализирует синлтон-объект для логирования.
func InitLogger() error {
	var err error
	Logger, err = util.NewLogger(&util.LoggingConfig{Logfile: "stdout", Level: "debug"})
	return err
}
