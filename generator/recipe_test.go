package generator_test

import (
	"bytes"
	_ "embed"
	"reflect"
	"testing"

	"github.com/goto/meteor/generator"
	"github.com/goto/meteor/plugins"
	"github.com/goto/meteor/registry"
	"github.com/goto/meteor/test/mocks"
	"github.com/stretchr/testify/assert"
)

var recipeVersions = [1]string{"v1beta1"}

func TestRecipe(t *testing.T) {
	if err := registry.Extractors.Register("test-extractor", func() plugins.Extractor {
		extr := mocks.NewExtractor()
		mockInfo := plugins.Info{
			Description: "Mock Extractor 1",
		}
		extr.On("Info").Return(mockInfo, nil).Once()
		return extr
	}); err != nil {
		t.Fatal(err)
	}

	if err := registry.Sinks.Register("test-sink", func() plugins.Syncer {
		mockSink := mocks.NewSink()
		mockInfo := plugins.Info{
			Description: "Mock Sink 1",
		}
		mockSink.On("Info").Return(mockInfo, nil).Once()
		return mockSink
	}); err != nil {
		t.Fatal(err)
	}

	if err := registry.Processors.Register("test-processor", func() plugins.Processor {
		mockProcessor := mocks.NewProcessor()
		mockInfo := plugins.Info{
			Description: "Mock Processor 1",
		}
		mockProcessor.On("Info").Return(mockInfo, nil).Once()
		return mockProcessor
	}); err != nil {
		t.Fatal(err)
	}

	type args struct {
		p generator.RecipeParams
	}
	tests := []struct {
		name    string
		args    args
		want    *generator.TemplateData
		wantErr bool
	}{
		{
			name: "success with minimal params",
			args: args{
				p: generator.RecipeParams{
					Name: "test-name",
				},
			},
			want: &generator.TemplateData{
				Name:    "test-name",
				Version: recipeVersions[len(recipeVersions)-1],
			},
			wantErr: false,
		},
		{
			name: "success with full params",
			args: args{
				p: generator.RecipeParams{
					Name:       "test-name",
					Source:     "test-extractor",
					Sinks:      []string{"test-sink"},
					Processors: []string{"test-processor"},
				},
			},
			want: &generator.TemplateData{
				Name:    "test-name",
				Version: recipeVersions[len(recipeVersions)-1],
				Source: struct {
					Name         string
					Scope        string
					SampleConfig string
				}{
					Name: "test-extractor",
				},
				Sinks: map[string]string{
					"test-sink": "",
				},
				Processors: map[string]string{
					"test-processor": "",
				},
			},
			wantErr: false,
		},
		{
			name: "error with invalid source",
			args: args{
				p: generator.RecipeParams{
					Name:   "test-name",
					Source: "invalid-source",
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error with invalid sinks",
			args: args{
				p: generator.RecipeParams{
					Name:  "test-name",
					Sinks: []string{"invalid-sink"},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error with invalid processors",
			args: args{
				p: generator.RecipeParams{
					Name:       "test-name",
					Processors: []string{"invalid-processor"},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generator.Recipe(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Recipe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Recipe() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecipeWriteTo(t *testing.T) {
	type args struct {
		p generator.RecipeParams
	}
	tests := []struct {
		name       string
		args       args
		wantWriter string
		wantErr    bool
	}{
		{
			name: "success with minimal params",
			args: args{
				p: generator.RecipeParams{
					Name: "test-name",
				},
			},
			wantWriter: `name: test-name
version: v1beta1
source:
  name: 
  config:     
    
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			if err := generator.RecipeWriteTo(tt.args.p, writer); (err != nil) != tt.wantErr {
				t.Errorf("RecipeWriteTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.wantWriter, writer.String())
		})
	}
}

func TestGetRecipeVersions(t *testing.T) {
	tests := []struct {
		name string
		want [1]string
	}{
		{
			name: "success",
			want: recipeVersions,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := generator.GetRecipeVersions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRecipeVersions() = %v, want %v", got, tt.want)
			}
		})
	}
}
