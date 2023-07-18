package metrics_test

import (
	"context"
	"testing"

	"github.com/goto/meteor/agent"
	"github.com/goto/meteor/config"
	"github.com/goto/meteor/metrics"
	"github.com/goto/meteor/recipe"
	"github.com/goto/salt/log"
	"github.com/stretchr/testify/assert"
)

func TestOtelMonitor_RecordRun(t *testing.T) {
	ctx := context.Background()

	t.Run("should create metrics with the correct name and value", func(t *testing.T) {
		recipe := recipe.Recipe{
			Name: "test-recipe",
			Source: recipe.PluginRecipe{
				Name: "mysql",
			},
		}
		duration := 100
		recordCount := 2

		monitor, done, err := metrics.NewOTLP(ctx, config.Config{
			OtelEnabled:       true,
			OtelCollectorAddr: "localhost:4317",
		}, log.NewLogrus(), "")
		defer done()

		monitor.RecordRun(ctx, agent.Run{Recipe: recipe, DurationInMs: duration, RecordCount: recordCount, Success: false})

		assert.Nil(t, err)
		assert.NotNil(t, monitor)
		assert.NotNil(t, done)
	})

}

func TestOtelMonitor_RecordPlugin(t *testing.T) {
	ctx := context.Background()
	t.Run("should create metrics with the correct name and value", func(t *testing.T) {
		recipe := recipe.Recipe{
			Name: "test-recipe",
			Source: recipe.PluginRecipe{
				Name: "bigquery",
			},
			Sinks: []recipe.PluginRecipe{
				{Name: "test-sink"},
			},
		}

		monitor, done, err := metrics.NewOTLP(ctx, config.Config{
			OtelEnabled:       true,
			OtelCollectorAddr: "localhost:4317",
		}, log.NewLogrus(), "")
		defer done()

		monitor.RecordPlugin(context.Background(),
			agent.PluginInfo{
				RecipeName: recipe.Name,
				PluginName: recipe.Sinks[0].Name,
				PluginType: "sink",
				Success:    true,
			})
		assert.Nil(t, err)
		assert.NotNil(t, monitor)
		assert.NotNil(t, done)
	})
}
