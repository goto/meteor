package config_test

import (
	"reflect"
	"testing"

	"github.com/goto/meteor/config"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	type args struct {
		configFile string
	}
	tests := []struct {
		name        string
		args        args
		expected    config.Config
		expectedErr string
	}{
		{
			name: "should return a config",
			args: args{
				configFile: "testdata/valid-config.yaml",
			},
			expected: config.Config{
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
			expected: config.Config{
				LogLevel:                    "info",
				StatsdEnabled:               false,
				StatsdHost:                  "localhost:8125",
				StatsdPrefix:                "meteor",
				MaxRetries:                  5,
				RetryInitialIntervalSeconds: 5,
			},
			expectedErr: "",
		},
		{
			name: "config invalid",
			args: args{
				configFile: "testdata/invalid-config.yaml",
			},
			expected:    config.Config{},
			expectedErr: "unable to load config to struct",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := config.Load(tt.args.configFile)
			if tt.expectedErr != "" {
				assert.ErrorContains(t, err, tt.expectedErr)
				return
			}

			assert.NoError(t, err)
			if !reflect.DeepEqual(actual, tt.expected) {
				t.Errorf("Load() = %v, want %v", actual, tt.expected)
			}
		})
	}
}