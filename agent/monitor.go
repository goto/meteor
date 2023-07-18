package agent

import (
	"context"
)

// Monitor is the interface for monitoring the agent.
type Monitor interface {
	RecordRun(ctx context.Context, run Run)
	RecordPlugin(ctx context.Context, pluginInfo PluginInfo)
}

// defaultMonitor is the default implementation of Monitor.
type defaultMonitor struct{}

func (*defaultMonitor) RecordRun(context.Context, Run) {
}

func (*defaultMonitor) RecordPlugin(context.Context, PluginInfo) {
}

func isNilMonitor(monitor []Monitor) bool {
	return len(monitor) == 0
}
