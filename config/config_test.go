package config_test

import (
	"reflect"
	"testing"

	"github.com/goto/meteor/config"
)

func TestLoad(t *testing.T) {
	type args struct {
		configFile string
	}
	tests := []struct {
		name    string
		args    args
		want    config.Config
		wantErr bool
	}{
		{
			name: "should return a config",
			args: args{
				configFile: "testdata/valid-config.yaml",
			},
			want: config.Config{
				LogLevel:                    "info",
				StatsdEnabled:               false,
				StatsdHost:                  "localhost:8125",
				StatsdPrefix:                "meteor",
				MaxRetries:                  5,
				RetryInitialIntervalSeconds: 5,
				StopOnSinkError:             false,
			},
		},
		{
			name: "config file not found",
			args: args{
				configFile: "not-found.yaml",
			},
			want: config.Config{
				LogLevel:                    "info",
				StatsdEnabled:               false,
				StatsdHost:                  "localhost:8125",
				StatsdPrefix:                "meteor",
				MaxRetries:                  5,
				RetryInitialIntervalSeconds: 5,
			},
			wantErr: false,
		},
		{
			name: "config invalid",
			args: args{
				configFile: "testdata/invalid-config.yaml",
			},
			want:    config.Config{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := config.Load(tt.args.configFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Load() = %v, want %v", got, tt.want)
			}
		})
	}
}
