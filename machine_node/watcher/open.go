package watcher

import (
	"context"

	"github.com/GDVFox/dflow/util"
)

// ActionWatcher объект синглтон для слежения за работоспособностью действий.
var ActionWatcher *Watcher

// StartWatcher инициализирует синглтон ActionWatcher и запускает его.
func StartWatcher(ctx context.Context, l *util.Logger, cfg *Config) error {
	ActionWatcher = newWatcher(l, cfg)
	return ActionWatcher.start(ctx)
}
