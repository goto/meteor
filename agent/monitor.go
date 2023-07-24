package agent

import (
	"context"
	"time"
)

type PluginInfo struct {
	RecipeName string
	PluginName string
	PluginType string
	Success    bool
	StartTime  time.Time
	BatchSize  int
}

// Monitor is the interface for monitoring the agent.
type Monitor interface {
	RecordRun(ctx context.Context, run Run)
	RecordPlugin(ctx context.Context, pluginInfo PluginInfo)
	RecordPluginRetryCount(ctx context.Context, pluginInfo PluginInfo)
}
