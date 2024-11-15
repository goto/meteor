//go:build plugins
// +build plugins

package maxcompute_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	// "github.com/aliyun/aliyun-odps-go-sdk/odps"

	v1beta2 "github.com/goto/meteor/models/gotocompany/assets/v1beta2"
	"github.com/goto/meteor/plugins"
	"github.com/goto/meteor/plugins/extractors/maxcompute"
	"github.com/goto/meteor/plugins/extractors/maxcompute/mocks"
	mocks2 "github.com/goto/meteor/test/mocks"
	"github.com/goto/meteor/test/utils"
	"github.com/goto/salt/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	projectID = "test-project-id"
)

// func mockClient(ctx context.Context, logger slog.Logger, config *maxcompute.Config) (maxcompute.Client, error) {
// 	client := mocks2.NewMaxComputeClient()
// 	return client, nil
// }

func TestInit(t *testing.T) {
	t.Run("should return error if config is invalid", func(t *testing.T) {
		extr := maxcompute.New(utils.Logger, maxcompute.CreateClient)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		err := extr.Init(ctx, plugins.Config{
			URNScope: "test-maxcompute",
			RawConfig: map[string]interface{}{
				"project_name": "",
			},
		})

		assert.ErrorContains(t, err, "project_name is required")
	})

	t.Run("should return no error", func(t *testing.T) {
		extr := maxcompute.New(utils.Logger, maxcompute.CreateClient)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		err := extr.Init(ctx, plugins.Config{
			URNScope: "test-maxcompute",
			RawConfig: map[string]interface{}{
				"project_name": projectID,
				"access_key": map[string]interface{}{
					"id":     "access_key_id",
					"secret": "access_key_secret",
				},
				"endpoint_project": "https://example.com/some-api",
			},
		})

		assert.NoError(t, err)
	})
}

func createClient(client *mocks.MaxComputeClient) func(ctx context.Context, logger log.Logger, config maxcompute.Config) (maxcompute.Client, error) {
	return func(ctx context.Context, logger log.Logger, config maxcompute.Config) (maxcompute.Client, error) {
		return client, nil
	}
}

func TestExtract(t *testing.T) {
	// schema1 := []*odps.Schema{
	// 	odps.NewSchema(nil, projectID, "my_schema"),
	// }

	// table1 := []*odps.Table{
	// 	odps.NewTable(nil, projectID, "my_schema", "my_dummy_table"),
	// 	odps.NewTable(nil, projectID, "my_schema", "my_dummy_table2"),
	// }

	runTest := func(t *testing.T, cfg plugins.Config, mockSetup func(mockClient *mocks.MaxComputeClient)) ([]*v1beta2.Asset, error) {
		mockClient := mocks.NewMaxComputeClient(t)
		if mockSetup != nil {
			mockSetup(mockClient)
		}
		extr := maxcompute.New(utils.Logger, createClient(mockClient))
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		err := extr.Init(ctx, cfg)
		if err != nil {
			return nil, err
		}

		emitter := mocks2.NewEmitter()
		err = extr.Extract(ctx, emitter.Push)

		actual := getAllData(emitter, t)
		return actual, err
	}

	// t.Run("should return no error", func(t *testing.T) {
	// 	actual, err := runTest(t, plugins.Config{
	// 		URNScope: "test-maxcompute",
	// 		RawConfig: map[string]interface{}{
	// 			"project_name": projectID,
	// 			"access_key": map[string]interface{}{
	// 				"id":     "access_key_id",
	// 				"secret": "access_key_secret",
	// 			},
	// 			"endpoint_project": "https://example.com/some-api",
	// 		},
	// 	}, func(mockClient *mocks.MaxComputeClient) {
	// 		mockClient.EXPECT().ListSchema(mock.Anything).Return(schema1, nil)
	// 		mockClient.EXPECT().ListTable(mock.Anything, "my_schema").Return(table1, nil)
	// 		mockClient.EXPECT().GetTable(mock.Anything, table1[0]).Return(table1[0], nil)
	// 		mockClient.EXPECT().GetTable(mock.Anything, table1[1]).Return(table1[1], nil)
	// 	})

	// 	assert.Nil(t, err)
	// 	assert.NotEmpty(t, actual)
	// 	fmt.Printf("********actual: %v\n", actual)
	// 	utils.AssertProtosWithJSONFile(t, "testdata/expected-assets.json", actual)
	// })

	// t.Run("should exclude schemas based on exclusion rules", func(t *testing.T) {
	// 	actual := runTest(t, plugins.Config{
	// 		URNScope: "test-maxcompute",
	// 		RawConfig: map[string]interface{}{
	// 			"project_name": projectID,
	// 			"access_key": map[string]interface{}{
	// 				"id":     "access_key_id",
	// 				"secret": "access_key_secret",
	// 			},
	// 			"endpoint_project": "https://example.com/some-api",
	// 			"exclude_schemas":  []string{"schema2"},
	// 		},
	// 	}, func(mockClient *mocks.MaxComputeClient) {
	// 		mockClient.EXPECT().ListSchema(mock.Anything).Return([]*odps.Schema{
	// 			{
	// 				Name: "schema1",
	// 				Tables: []*odps.Table{
	// 					{Name: "table1", Comment: "Test table 1"},
	// 					{Name: "table2", Comment: "Test table 2"},
	// 				},
	// 			},
	// 			{
	// 				Name: "schema2",
	// 				Tables: []*odps.Table{
	// 					{Name: "table3", Comment: "Test table 3"},
	// 				},
	// 			},
	// 		}, nil)
	// 	})

	// 	utils.AssertProtosWithJSONFile(t, "testdata/expected-assets-excluded-schema.json", actual)
	// })

	// t.Run("should exclude tables based on exclusion rules", func(t *testing.T) {
	// 	actual := runTest(t, plugins.Config{
	// 		URNScope: "test-maxcompute",
	// 		RawConfig: map[string]interface{}{
	// 			"project_name": projectID,
	// 			"access_key": map[string]interface{}{
	// 				"id":     "access_key_id",
	// 				"secret": "access_key_secret",
	// 			},
	// 			"endpoint_project": "https://example.com/some-api",
	// 			"exclude_tables":   []string{"schema1.table2"},
	// 		},
	// 	}, func(mockClient *mocks.MaxComputeClient) {
	// 		mockClient.EXPECT().ListSchema(mock.Anything).Return([]*odps.Schema{
	// 			{
	// 				Name: "schema1",
	// 				Tables: []*odps.Table{
	// 					{Name: "table1", Comment: "Test table 1"},
	// 					{Name: "table2", Comment: "Test table 2"},
	// 				},
	// 			},
	// 		}, nil)
	// 	})

	// 	utils.AssertProtosWithJSONFile(t, "testdata/expected-assets-excluded-table.json", actual)
	// })

	t.Run("should return error if ListSchema fails", func(t *testing.T) {
		actual, err := runTest(t, plugins.Config{
			URNScope: "test-maxcompute",
			RawConfig: map[string]interface{}{
				"project_name": projectID,
				"access_key": map[string]interface{}{
					"id":     "access_key_id",
					"secret": "access_key_secret",
				},
				"endpoint_project": "https://example.com/some-api",
			},
		}, func(mockClient *mocks.MaxComputeClient) {
			mockClient.EXPECT().ListSchema(mock.Anything).Return(nil, fmt.Errorf("ListSchema fails"))
		})
		assert.ErrorContains(t, err, "ListSchema fails")
		assert.Nil(t, actual)
	})
}

func getAllData(emitter *mocks2.Emitter, t *testing.T) []*v1beta2.Asset {
	actual := emitter.GetAllData()

	// the emulator appending 1 random dataset
	// we can't assert it, so we remove it from the list
	// if len(actual) > 0 {
	// 	actual = actual[:len(actual)-1]
	// }

	// the emulator returning dynamic timestamps
	// replace them with static ones
	for _, asset := range actual {
		replaceWithStaticTimestamp(t, asset)
	}
	return actual
}

func replaceWithStaticTimestamp(t *testing.T, asset *v1beta2.Asset) {
	b := new(v1beta2.Table)
	err := asset.Data.UnmarshalTo(b)
	assert.NoError(t, err)

	time, err := time.Parse(time.RFC3339, "2023-06-13T03:46:12.372974Z")
	assert.NoError(t, err)
	b.CreateTime = timestamppb.New(time)
	b.UpdateTime = timestamppb.New(time)

	asset.Data, err = anypb.New(b)
	assert.NoError(t, err)
}
