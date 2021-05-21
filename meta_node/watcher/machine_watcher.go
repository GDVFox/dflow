package watcher

import (
	"context"
	"fmt"
	"time"

	"github.com/GDVFox/dflow/meta_node/planner"
	"github.com/GDVFox/dflow/util"
	"github.com/GDVFox/dflow/util/message"
	"github.com/pkg/errors"
)

// Возможные ошибки.
var (
	ErrNoHost             = errors.New("unknown machine host")
	ErrPlanAlreadyStarted = errors.New("plan already started")
)

// MachineWatcherConfig набор настроек для Watcher
type MachineWatcherConfig struct {
	PingFrequency util.Duration    `yaml:"ping-freq"`
	Machines      []*MachineConfig `yaml:"machines"`
}

// NewMachineWatcherConfig возвращает MachineWatcherConfig с настройками по умолчанию.
func NewMachineWatcherConfig() *MachineWatcherConfig {
	return &MachineWatcherConfig{
		PingFrequency: util.Duration(5 * time.Second),
		Machines:      make([]*MachineConfig, 0),
	}
}

// MachineWatcher структура для контроля запущенных действий.
type MachineWatcher struct {
	machines map[string]*Machine

	cfg    *MachineWatcherConfig
	logger *util.Logger
}

// NewWatcher создает новый объект watcher
func newMachineWatcher(l *util.Logger, cfg *MachineWatcherConfig) (*MachineWatcher, error) {
	machines := make(map[string]*Machine, len(cfg.Machines))
	for _, machineConfig := range cfg.Machines {
		if _, ok := machines[machineConfig.Host]; ok {
			return nil, fmt.Errorf("duplicate host: %s", machineConfig.Host)
		}
		machines[machineConfig.Host] = NewMachine(l, machineConfig)
	}

	return &MachineWatcher{
		machines: machines,
		cfg:      cfg,
		logger:   l,
	}, nil
}

func (w *MachineWatcher) sendRunAction(ctx context.Context, schemeName string, node *planner.NodePlan) error {
	w.logger.Debugf("machine_watcher: starting action %s from plan %s on %s:%d", node.Action, schemeName, node.Host, node.Port)

	machine, ok := w.machines[node.Host]
	if !ok {
		return ErrNoHost
	}
	return machine.SendRunAction(ctx, schemeName, node)
}

func (w *MachineWatcher) sendStopAction(ctx context.Context, schemeName string, node *planner.NodePlan) error {
	w.logger.Debugf("machine_watcher: stopping action %s from plan %s on %s:%d", node.Action, schemeName, node.Host, node.Port)

	machine, ok := w.machines[node.Host]
	if !ok {
		return ErrNoHost
	}
	return machine.SendStopAction(ctx, schemeName, node)
}

func (w *MachineWatcher) sendChangeOut(ctx context.Context, schemeName, oldOut, newOut string, node *planner.NodePlan) error {
	w.logger.Debugf("machine_watcher: action %s from plan %s changing out (%s -> %s)", node.Action, schemeName, oldOut, newOut)

	machine, ok := w.machines[node.Host]
	if !ok {
		return ErrNoHost
	}
	return machine.SendChangeOut(ctx, schemeName, node.Name, oldOut, newOut)
}

func (w *MachineWatcher) pingMachines() map[string]*message.RuntimeTelemetry {
	w.logger.Debugf("machine_watcher: started ping machines")

	runtimes := make(map[string]*message.RuntimeTelemetry)
	for machineHost, machine := range w.machines {
		telemetry, err := machine.Ping()
		if err != nil {
			w.logger.Warnf("machine_watcher: ping for machine '%s' failed: %v", machineHost, err)
			continue
		}

		for _, metrics := range telemetry {
			runtimeName := buildRuntimeName(metrics.SchemeName, metrics.ActionName)
			runtimes[runtimeName] = metrics
		}
	}

	w.logger.Debugf("machine_watcher: ping machines done")
	return runtimes
}

func buildRuntimeName(schemeName, actionName string) string {
	return schemeName + "_" + actionName
}
