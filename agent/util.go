package agent

import (
	"github.com/goto/meteor/plugins"
	"github.com/goto/meteor/recipe"
)

func recipeToPluginConfig(pr recipe.PluginRecipe, oe bool) plugins.Config {
	return plugins.Config{
		URNScope:    pr.Scope,
		RawConfig:   pr.Config,
		OtelEnabled: oe,
	}
}
