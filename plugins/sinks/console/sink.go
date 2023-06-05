package console

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/goto/meteor/models"
	assetsv1beta2 "github.com/goto/meteor/models/gotocompany/assets/v1beta2"
	"github.com/goto/meteor/plugins"
	"github.com/goto/meteor/registry"
	"github.com/goto/salt/log"
)

//go:embed README.md
var summary string

var info = plugins.Info{
	Description:  "Log to standard output",
	Summary:      summary,
	Tags:         []string{"log", "sink"},
	SampleConfig: "",
}

type Sink struct {
	plugins.BasePlugin
	logger log.Logger
}

func New(logger log.Logger) plugins.Syncer {
	s := &Sink{
		logger: logger,
	}
	s.BasePlugin = plugins.NewBasePlugin(info, nil)

	return s
}

func (s *Sink) Init(ctx context.Context, config plugins.Config) error {
	return s.BasePlugin.Init(ctx, config)
}

func (s *Sink) Sink(_ context.Context, batch []models.Record) error {
	for _, record := range batch {
		if err := s.process(record.Data()); err != nil {
			return err
		}
	}
	return nil
}

func (*Sink) Close() error { return nil }

func (*Sink) process(asset *assetsv1beta2.Asset) error {
	jsonBytes, err := models.ToJSON(asset)
	if err != nil {
		return err
	}

	fmt.Println(string(jsonBytes))

	return nil
}

func init() {
	if err := registry.Sinks.Register("console", func() plugins.Syncer {
		return New(plugins.GetLog())
	}); err != nil {
		panic(err)
	}
}
