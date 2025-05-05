//go:build plugins
// +build plugins

package maxcompute_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aliyun/aliyun-odps-go-sdk/odps"
	"github.com/aliyun/aliyun-odps-go-sdk/odps/common"
	"github.com/aliyun/aliyun-odps-go-sdk/odps/datatype"
	"github.com/aliyun/aliyun-odps-go-sdk/odps/tableschema"
	v1beta2 "github.com/goto/meteor/models/gotocompany/assets/v1beta2"
	"github.com/goto/meteor/plugins"
	"github.com/goto/meteor/plugins/extractors/maxcompute"
	"github.com/goto/meteor/plugins/extractors/maxcompute/config"
	"github.com/goto/meteor/plugins/extractors/maxcompute/mocks"
	mocks2 "github.com/goto/meteor/test/mocks"
	"github.com/goto/meteor/test/utils"
	"github.com/goto/salt/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	projectID = "test-project-id"
)

func TestInit(t *testing.T) {

	mockClient := mocks.NewMaxComputeClient(t)
	t.Run("should return error if config is invalid", func(t *testing.T) {
		extr := maxcompute.New(utils.Logger, createClient(mockClient), nil)
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
		extr := maxcompute.New(utils.Logger, createClient(mockClient), nil)
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

func createClient(client *mocks.MaxComputeClient) func(ctx context.Context, logger log.Logger, config config.Config) (maxcompute.Client, error) {
	return func(ctx context.Context, logger log.Logger, config config.Config) (maxcompute.Client, error) {
		return client, nil
	}
}

func TestExtract(t *testing.T) {
	schema1 := []*odps.Schema{
		odps.NewSchema(nil, projectID, "my_schema"),
	}

	table1 := []*odps.Table{
		odps.NewTable(nil, projectID, "my_schema", "dummy_table"),
		odps.NewTable(nil, projectID, "my_schema", "new_table"),
	}

	c1 := tableschema.Column{
		Name:    "id",
		Type:    datatype.BigIntType,
		NotNull: false,
	}

	c2 := tableschema.Column{
		Name: "name",
		Type: datatype.StructType{
			Fields: []datatype.StructFieldType{
				{
					Name: "first_name",
					Type: datatype.StringType,
				},
				{
					Name: "last_name",
					Type: datatype.StringType,
				},
			},
		},
		NotNull: false,
	}

	// Schema for dummy_table
	dummyTableSchemaBuilder := tableschema.NewSchemaBuilder()
	dummyTableSchemaBuilder.Name("dummy_table").
		Columns(c1, c2)
	dummyTableSchema := dummyTableSchemaBuilder.Build()
	dummyTableSchema.ViewText = "SELECT id, name, user_info\nFROM test-project-id.default.my_dummy_table"
	dummyCreateTime, err := time.Parse(time.RFC3339, "2024-11-14T06:41:35Z")
	if err != nil {
		t.Fatalf("failed to parse create time for dummy_table: %v", err)
	}
	dummyTableSchema.CreateTime = common.GMTTime(dummyCreateTime)
	dummyTableSchema.LastModifiedTime = common.GMTTime(dummyCreateTime)
	dummyTableSchema.Comment = "dummy table description"

	c3 := tableschema.Column{
		Name:    "user_id",
		Type:    datatype.BigIntType,
		Comment: "Unique identifier for users",
		NotNull: true,
	}

	c4 := tableschema.Column{
		Name:    "email",
		Type:    datatype.StringType,
		NotNull: true,
		Comment: "User email address",
	}

	// Schema for new_table
	newTableSchemaBuilder := tableschema.NewSchemaBuilder()
	newTableSchemaBuilder.Name("new_table").
		Columns(c3, c4)
	newTableSchema := newTableSchemaBuilder.Build()
	newTableSchema.ViewText = "SELECT user_id, email FROM test-project-id.my_schema.new_table"
	newCreateTime, err := time.Parse(time.RFC3339, "2024-11-18T08:00:00Z")
	if err != nil {
		t.Fatalf("failed to parse create time for new_table: %v", err)
	}
	newTableSchema.CreateTime = common.GMTTime(newCreateTime)
	newTableSchema.LastModifiedTime = common.GMTTime(newCreateTime)

	// Schema mapping
	schemaMapping := map[string]*tableschema.TableSchema{
		"dummy_table": &dummyTableSchema,
		"new_table":   &newTableSchema,
	}

	runTest := func(t *testing.T, cfg plugins.Config, mockSetup func(mockClient *mocks.MaxComputeClient), randomizer func(seed int64) func(int64) int64) ([]*v1beta2.Asset, error) {
		mockClient := mocks.NewMaxComputeClient(t)
		if mockSetup != nil {
			mockSetup(mockClient)
		}
		extr := maxcompute.New(utils.Logger, createClient(mockClient), randomizer)
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		err := extr.Init(ctx, cfg)
		if err != nil {
			return nil, err
		}

		emitter := mocks2.NewEmitter()
		err = extr.Extract(ctx, emitter.Push)

		actual := emitter.GetAllData()
		return actual, err
	}

	t.Run("should return no error without lineage", func(t *testing.T) {
		actual, err := runTest(t, plugins.Config{
			URNScope: "test-maxcompute",
			RawConfig: map[string]interface{}{
				"project_name": projectID,
				"access_key": map[string]interface{}{
					"id":     "access_key_id",
					"secret": "access_key_secret",
				},
				"endpoint_project":   "https://example.com/some-api",
				"max_preview_rows":   3,
				"mix_values":         false,
				"build_view_lineage": false,
			},
		}, func(mockClient *mocks.MaxComputeClient) {
			mockClient.EXPECT().ListSchema(mock.Anything).Return(schema1, nil)

			mockClient.EXPECT().ListTable(mock.Anything, "my_schema").Return(table1, nil)
			mockClient.EXPECT().GetTableSchema(mock.Anything, table1[0]).Return("VIRTUAL_VIEW", schemaMapping[table1[0].Name()], nil)
			mockClient.EXPECT().GetTablePreview(mock.Anything, "", table1[0], 3).Return(nil, nil, nil)

			mockClient.EXPECT().ListTable(mock.Anything, "my_schema").Return(table1[1:], nil)
			mockClient.EXPECT().GetTableSchema(mock.Anything, table1[1]).Return("MANAGED_TABLE", schemaMapping[table1[1].Name()], nil)
			mockClient.EXPECT().GetTablePreview(mock.Anything, "", table1[1], 3).Return(
				[]string{"user_id", "email"},
				&structpb.ListValue{
					Values: []*structpb.Value{
						structpb.NewListValue(&structpb.ListValue{
							Values: []*structpb.Value{
								structpb.NewStringValue("1"),
								structpb.NewStringValue("user1@example.com"),
							},
						}),
						structpb.NewListValue(&structpb.ListValue{
							Values: []*structpb.Value{
								structpb.NewStringValue("2"),
								structpb.NewStringValue("user2@example.com"),
							},
						}),
					},
				},
				nil,
			)
			mockClient.EXPECT().GetMaskingPolicies(mock.Anything).Return(map[string][]string{
				"user_id": {"policyTag1"},
				"email":   {"policyTag2", "policyTag3"},
			}, nil)
		}, nil)

		assert.Nil(t, err)
		assert.NotEmpty(t, actual)
		utils.AssertProtosWithJSONFile(t, "testdata/expected-assets.json", actual)
	})

	t.Run("should return no error with lineage", func(t *testing.T) {
		actual, err := runTest(t, plugins.Config{
			URNScope: "test-maxcompute",
			RawConfig: map[string]interface{}{
				"project_name": projectID,
				"access_key": map[string]interface{}{
					"id":     "access_key_id",
					"secret": "access_key_secret",
				},
				"endpoint_project":   "https://example.com/some-api",
				"max_preview_rows":   3,
				"mix_values":         false,
				"build_view_lineage": true,
			},
		}, func(mockClient *mocks.MaxComputeClient) {
			mockClient.EXPECT().ListSchema(mock.Anything).Return(schema1, nil)

			mockClient.EXPECT().ListTable(mock.Anything, "my_schema").Return(table1, nil)
			mockClient.EXPECT().GetTableSchema(mock.Anything, table1[0]).Return("VIRTUAL_VIEW", schemaMapping[table1[0].Name()], nil)
			mockClient.EXPECT().GetTablePreview(mock.Anything, "", table1[0], 3).Return(nil, nil, nil)

			mockClient.EXPECT().ListTable(mock.Anything, "my_schema").Return(table1[1:], nil)
			mockClient.EXPECT().GetTableSchema(mock.Anything, table1[1]).Return("MANAGED_TABLE", schemaMapping[table1[1].Name()], nil)
			mockClient.EXPECT().GetTablePreview(mock.Anything, "", table1[1], 3).Return(
				[]string{"user_id", "email"},
				&structpb.ListValue{
					Values: []*structpb.Value{
						structpb.NewListValue(&structpb.ListValue{
							Values: []*structpb.Value{
								structpb.NewStringValue("1"),
								structpb.NewStringValue("user1@example.com"),
							},
						}),
						structpb.NewListValue(&structpb.ListValue{
							Values: []*structpb.Value{
								structpb.NewStringValue("2"),
								structpb.NewStringValue("user2@example.com"),
							},
						}),
					},
				},
				nil,
			)
			mockClient.EXPECT().GetMaskingPolicies(mock.Anything).Return(map[string][]string{}, nil)
		}, nil)

		assert.Nil(t, err)
		assert.NotEmpty(t, actual)
		utils.AssertProtosWithJSONFile(t, "testdata/expected-assets-with-view-lineage.json", actual)
	})

	t.Run("should exclude tables based on exclusion rules", func(t *testing.T) {
		actual, err := runTest(t, plugins.Config{
			URNScope: "test-maxcompute",
			RawConfig: map[string]interface{}{
				"project_name": projectID,
				"access_key": map[string]interface{}{
					"id":     "access_key_id",
					"secret": "access_key_secret",
				},
				"endpoint_project": "https://example.com/some-api",
				"exclude": map[string]interface{}{
					"tables": []string{"my_schema.dummy_table"},
				},
			},
		}, func(mockClient *mocks.MaxComputeClient) {
			mockClient.EXPECT().ListSchema(mock.Anything).Return(schema1, nil)
			mockClient.EXPECT().ListTable(mock.Anything, "my_schema").Return(table1[1:], nil)
			mockClient.EXPECT().GetTableSchema(mock.Anything, table1[1]).Return("MANAGED_TABLE", schemaMapping[table1[1].Name()], nil)
			mockClient.EXPECT().GetMaskingPolicies(mock.Anything).Return(map[string][]string{}, nil)
		}, nil)

		assert.Nil(t, err)
		assert.NotEmpty(t, actual)
		utils.AssertProtosWithJSONFile(t, "testdata/expected-assets-with-table-exclusion.json", actual)
	})

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
		}, nil)
		assert.ErrorContains(t, err, "ListSchema fails")
		assert.Nil(t, actual)
	})
}
