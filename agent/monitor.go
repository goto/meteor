package agent

import (
	"context"
)

// Monitor is the interface for monitoring the agent.
type Monitor interface {
	RecordRun(ctx context.Context, run Run)
	RecordPlugin(ctx context.Context, pluginInfo PluginInfo)
}
