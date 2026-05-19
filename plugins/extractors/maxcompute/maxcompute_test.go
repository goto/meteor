//go:build plugins
// +build plugins

package maxcompute_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

const (
	projectID = "test-project-id"
)

func createClient(client *mocks.MaxComputeClient) func(ctx context.Context, logger log.Logger, config config.Config) (maxcompute.Client, error) {
	return func(ctx context.Context, logger log.Logger, config config.Config) (maxcompute.Client, error) {
		return client, nil
	}
}

// newGroupMappingServer starts a test HTTP server that responds to the Shield
// groups endpoint. Pass nil groups for an empty response (no PDG enrichment).
func newGroupMappingServer(t *testing.T, groups []map[string]interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		body := map[string]interface{}{"groups": groups}
		if groups == nil {
			body["groups"] = []interface{}{}
		}
		if err := json.NewEncoder(w).Encode(body); err != nil {
			http.Error(w, "encode error", http.StatusInternalServerError)
		}
	}))
}

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

	t.Run("should return error if access_key id is empty", func(t *testing.T) {
		extr := maxcompute.New(utils.Logger, createClient(mockClient), nil)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		err := extr.Init(ctx, plugins.Config{
			URNScope: "test-maxcompute",
			RawConfig: map[string]interface{}{
				"project_name": projectID,
				"access_key": map[string]interface{}{
					"id":     "",
					"secret": "access_key_secret",
				},
				"endpoint_project": "https://example.com/some-api",
			},
		})

		assert.ErrorContains(t, err, "access_key is required")
	})

	t.Run("should return error if access_key secret is empty", func(t *testing.T) {
		extr := maxcompute.New(utils.Logger, createClient(mockClient), nil)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		err := extr.Init(ctx, plugins.Config{
			URNScope: "test-maxcompute",
			RawConfig: map[string]interface{}{
				"project_name": projectID,
				"access_key": map[string]interface{}{
					"id":     "access_key_id",
					"secret": "",
				},
				"endpoint_project": "https://example.com/some-api",
			},
		})

		assert.ErrorContains(t, err, "access_key is required")
	})

	t.Run("should return error if endpoint_project is empty", func(t *testing.T) {
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
				"endpoint_project": "",
			},
		})

		assert.ErrorContains(t, err, "endpoint_project is required")
	})

	t.Run("should return no error", func(t *testing.T) {
		shieldServer := newGroupMappingServer(t, nil)
		defer shieldServer.Close()

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
				"shield_host":      shieldServer.URL,
			},
		})

		assert.NoError(t, err)
	})
}

func TestExtract(t *testing.T) {
	schema1 := []*odps.Schema{
		odps.NewSchema(nil, projectID, "my_schema"),
	}

	table1 := []*odps.Table{
		odps.NewTable(nil, projectID, "my_schema", "dummy_table"),
		odps.NewTable(nil, projectID, "my_schema", "new_table"),
	}

	table2 := []*odps.Table{
		odps.NewTable(nil, projectID, "my_schema", "dummy_table"),
		odps.NewTable(nil, projectID, "my_schema", "new_table"),
		odps.NewTable(nil, projectID, "my_schema", "table_lifecycle_3"),
		odps.NewTable(nil, projectID, "my_schema", "table_lifecycle_8"),
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
	dummyTableSchema.Comment = `{"description": "dummy table description", "data_owner_team": "dummy_team"}`
	dummyTableSchema.PartitionColumns = []tableschema.Column{
		{
			Name: "data_date",
			Type: datatype.DateType,
		},
	}

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
	newTableSchema.TableName = "new_table"
	newTableSchema.ViewText = "SELECT user_id, email FROM test-project-id.my_schema.new_table"
	newCreateTime, err := time.Parse(time.RFC3339, "2024-11-18T08:00:00Z")
	if err != nil {
		t.Fatalf("failed to parse create time for new_table: %v", err)
	}
	newTableSchema.CreateTime = common.GMTTime(newCreateTime)
	newTableSchema.LastModifiedTime = common.GMTTime(newCreateTime)

	// Schema for table_lifecycle_3
	tableLifecycle3Schema := newTableSchema // copy
	tableLifecycle3Schema.TableName = "table_lifecycle_3"
	tableLifecycle3Schema.ViewText = "SELECT user_id, email FROM test-project-id.my_schema.table_lifecycle_3"
	tableLifecycle3Schema.Lifecycle = 3
	tableLifecycle3Schema.RecordNum = 100

	// Schema for table_lifecycle_8
	tableLifecycle8Schema := newTableSchema // copy
	tableLifecycle8Schema.TableName = "table_lifecycle_8"
	tableLifecycle8Schema.ViewText = "SELECT user_id, email FROM test-project-id.my_schema.table_lifecycle_8"
	tableLifecycle8Schema.Lifecycle = 8
	tableLifecycle8Schema.RecordNum = 200

	// Schema mapping
	schemaMapping := map[string]*tableschema.TableSchema{
		"dummy_table":       &dummyTableSchema,
		"new_table":         &newTableSchema,
		"table_lifecycle_3": &tableLifecycle3Schema,
		"table_lifecycle_8": &tableLifecycle8Schema,
	}

	runTest := func(t *testing.T, cfg plugins.Config, mockSetup func(mockClient *mocks.MaxComputeClient), randomizer func(seed int64) func(int64) int64) ([]*v1beta2.Asset, error) {
		// Inject a default empty-groups shield server unless the caller already provided shield_host.
		if _, hasShieldHost := cfg.RawConfig["shield_host"]; !hasShieldHost {
			shieldServer := newGroupMappingServer(t, nil)
			t.Cleanup(shieldServer.Close)
			cfg.RawConfig["shield_host"] = shieldServer.URL
		}

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

	t.Run("should return empty group mapping when shield_host is not configured", func(t *testing.T) {
		mockClient := mocks.NewMaxComputeClient(t)
		mockClient.EXPECT().ListSchema(mock.Anything).Return(schema1, nil)
		mockClient.EXPECT().ListTable(mock.Anything, "my_schema").Return(table1[:1], nil)
		mockClient.EXPECT().GetTableSchema(mock.Anything, table1[0]).Return("VIRTUAL_VIEW", schemaMapping[table1[0].Name()], nil)
		mockClient.EXPECT().GetMaskingPolicies(mock.Anything).Return(map[string][]string{}, nil)

		extr := maxcompute.New(utils.Logger, createClient(mockClient), nil)
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
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
		require.NoError(t, err)

		emitter := mocks2.NewEmitter()
		require.NoError(t, extr.Extract(ctx, emitter.Push))

		actual := emitter.GetAllData()
		require.Len(t, actual, 1)

		tableData := &v1beta2.Table{}
		require.NoError(t, proto.Unmarshal(actual[0].GetData().GetValue(), tableData))
		assert.NotContains(t, tableData.GetLabels(), "pdg")
	})

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
			mockClient.EXPECT().GetTablePreview(mock.Anything, "", table1[0], 3).Return(nil, nil, fmt.Errorf("failed to get table preview"))

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
					"tables":              []string{"my_schema.dummy_table"},
					"min_table_lifecycle": 8,
				},
			},
		}, func(mockClient *mocks.MaxComputeClient) {
			mockClient.EXPECT().ListSchema(mock.Anything).Return(schema1, nil)
			mockClient.EXPECT().ListTable(mock.Anything, "my_schema").Return(table2, nil)
			mockClient.EXPECT().GetTableSchema(mock.Anything, table2[1]).Return("MANAGED_TABLE", schemaMapping[table2[1].Name()], nil)
			mockClient.EXPECT().GetTableSchema(mock.Anything, table2[2]).Return("MANAGED_TABLE", schemaMapping[table2[2].Name()], nil)
			mockClient.EXPECT().GetTableSchema(mock.Anything, table2[3]).Return("MANAGED_TABLE", schemaMapping[table2[3].Name()], nil)
			mockClient.EXPECT().GetMaskingPolicies(mock.Anything).Return(map[string][]string{}, nil)
		}, nil)

		assert.Nil(t, err)
		assert.NotEmpty(t, actual)
		utils.AssertProtosWithJSONFile(t, "testdata/expected-assets-with-table-exclusion.json", actual)
	})

	t.Run("should return no error if GetTablePreview fails", func(t *testing.T) {
		actual, err := runTest(t, plugins.Config{
			URNScope: "test-maxcompute",
			RawConfig: map[string]interface{}{
				"project_name": projectID,
				"access_key": map[string]interface{}{
					"id":     "access_key_id",
					"secret": "access_key_secret",
				},
				"endpoint_project": "https://example.com/some-api",
				"max_preview_rows": 3,
			},
		}, func(mockClient *mocks.MaxComputeClient) {
			mockClient.EXPECT().ListSchema(mock.Anything).Return(schema1, nil)
			mockClient.EXPECT().ListTable(mock.Anything, "my_schema").Return(table1[1:2], nil)
			mockClient.EXPECT().GetTableSchema(mock.Anything, table1[1]).Return("MANAGED_TABLE", schemaMapping[table1[1].Name()], nil)
			mockClient.EXPECT().GetTablePreview(mock.Anything, "", table1[1], 3).Return(nil, nil, fmt.Errorf("preview failed"))
			mockClient.EXPECT().GetMaskingPolicies(mock.Anything).Return(map[string][]string{}, nil)
		}, nil)

		require.NoError(t, err)
		require.NotEmpty(t, actual)
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

	t.Run("should extract partition column name from generate expression", func(t *testing.T) {
		autoPartitionTable := odps.NewTable(nil, projectID, "my_schema", "auto_partition_table")

		autoPartitionSchemaBuilder := tableschema.NewSchemaBuilder()
		autoPartitionSchemaBuilder.Name("auto_partition_table").Columns(c3, c4)
		autoPartitionSchema := autoPartitionSchemaBuilder.Build()
		autoPartitionSchema.TableName = "auto_partition_table"
		autoPartitionSchema.CreateTime = common.GMTTime(newCreateTime)
		autoPartitionSchema.LastModifiedTime = common.GMTTime(newCreateTime)
		autoPartitionSchema.PartitionColumns = []tableschema.Column{
			{
				Name:               "_partition_value",
				Type:               datatype.StringType,
				GenerateExpression: tableschema.NewTruncTime("start_date", tableschema.DAY),
			},
		}

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
			mockClient.EXPECT().ListSchema(mock.Anything).Return(schema1, nil)
			mockClient.EXPECT().ListTable(mock.Anything, "my_schema").Return([]*odps.Table{autoPartitionTable}, nil)
			mockClient.EXPECT().GetTableSchema(mock.Anything, autoPartitionTable).Return("MANAGED_TABLE", &autoPartitionSchema, nil)
			mockClient.EXPECT().GetMaskingPolicies(mock.Anything).Return(map[string][]string{}, nil)
		}, nil)

		assert.Nil(t, err)
		assert.NotEmpty(t, actual)
		utils.AssertProtosWithJSONFile(t, "testdata/expected-assets-auto-partition.json", actual)
	})

	t.Run("multi-schema: all assets from schema A emitted before schema B is processed", func(t *testing.T) {
		schemaA := odps.NewSchema(nil, projectID, "schema_a")
		schemaB := odps.NewSchema(nil, projectID, "schema_b")
		tableA := odps.NewTable(nil, projectID, "schema_a", "table_a")
		tableB := odps.NewTable(nil, projectID, "schema_b", "table_b")

		tableASchemaBuilder := tableschema.NewSchemaBuilder()
		tableASchemaBuilder.Name("table_a").Columns(c3)
		tableASchema := tableASchemaBuilder.Build()
		tableASchema.TableName = "table_a"

		tableBSchemaBuilder := tableschema.NewSchemaBuilder()
		tableBSchemaBuilder.Name("table_b").Columns(c4)
		tableBSchema := tableBSchemaBuilder.Build()
		tableBSchema.TableName = "table_b"

		actual, err := runTest(t, plugins.Config{
			URNScope: "test-maxcompute",
			RawConfig: map[string]interface{}{
				"project_name": projectID,
				"access_key": map[string]interface{}{
					"id":     "access_key_id",
					"secret": "access_key_secret",
				},
				"endpoint_project": "https://example.com/some-api",
				"concurrency":      2,
			},
		}, func(mockClient *mocks.MaxComputeClient) {
			mockClient.EXPECT().ListSchema(mock.Anything).Return([]*odps.Schema{schemaA, schemaB}, nil)

			mockClient.EXPECT().ListTable(mock.Anything, "schema_a").Return([]*odps.Table{tableA}, nil)
			mockClient.EXPECT().GetTableSchema(mock.Anything, tableA).Return("MANAGED_TABLE", &tableASchema, nil)
			mockClient.EXPECT().GetMaskingPolicies(mock.Anything).Return(map[string][]string{}, nil).Once()

			mockClient.EXPECT().ListTable(mock.Anything, "schema_b").Return([]*odps.Table{tableB}, nil)
			mockClient.EXPECT().GetTableSchema(mock.Anything, tableB).Return("MANAGED_TABLE", &tableBSchema, nil)
			mockClient.EXPECT().GetMaskingPolicies(mock.Anything).Return(map[string][]string{}, nil).Once()
		}, nil)

		assert.NoError(t, err)
		assert.Len(t, actual, 2)

		urns := []string{actual[0].Urn, actual[1].Urn}
		assert.Contains(t, urns, plugins.MaxComputeURN(projectID, "schema_a", "table_a"))
		assert.Contains(t, urns, plugins.MaxComputeURN(projectID, "schema_b", "table_b"))
	})

	t.Run("should skip table with schema error and continue processing remaining schemas", func(t *testing.T) {
		schemaA := odps.NewSchema(nil, projectID, "schema_a")
		schemaB := odps.NewSchema(nil, projectID, "schema_b")
		tableA := odps.NewTable(nil, projectID, "schema_a", "table_a")
		tableB := odps.NewTable(nil, projectID, "schema_b", "table_b")

		tableBSchemaBuilder := tableschema.NewSchemaBuilder()
		tableBSchemaBuilder.Name("table_b").Columns(c3)
		tableBSchema := tableBSchemaBuilder.Build()
		tableBSchema.TableName = "table_b"

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
			mockClient.EXPECT().ListSchema(mock.Anything).Return([]*odps.Schema{schemaA, schemaB}, nil)

			mockClient.EXPECT().ListTable(mock.Anything, "schema_a").Return([]*odps.Table{tableA}, nil)
			mockClient.EXPECT().GetTableSchema(mock.Anything, tableA).Return("", nil, fmt.Errorf("schema_a table error"))

			mockClient.EXPECT().ListTable(mock.Anything, "schema_b").Return([]*odps.Table{tableB}, nil)
			mockClient.EXPECT().GetTableSchema(mock.Anything, tableB).Return("MANAGED_TABLE", &tableBSchema, nil)
			mockClient.EXPECT().GetMaskingPolicies(mock.Anything).Return(map[string][]string{}, nil)
		}, nil)

		assert.NoError(t, err)
		require.Len(t, actual, 1)
		assert.Equal(t, plugins.MaxComputeURN(projectID, "schema_b", "table_b"), actual[0].Urn)
	})

	t.Run("should enrich table labels with pdg when group mapping matches data_owner_team", func(t *testing.T) {
		shieldServer := newGroupMappingServer(t, []map[string]interface{}{
			{
				"slug": "dummy_team",
				"metadata": map[string]interface{}{
					"product-group-name": "Dummy PDG",
				},
			},
		})
		defer shieldServer.Close()

		actual, err := runTest(t, plugins.Config{
			URNScope: "test-maxcompute",
			RawConfig: map[string]interface{}{
				"project_name": projectID,
				"access_key": map[string]interface{}{
					"id":     "access_key_id",
					"secret": "access_key_secret",
				},
				"endpoint_project": "https://example.com/some-api",
				"shield_host":      shieldServer.URL,
			},
		}, func(mockClient *mocks.MaxComputeClient) {
			mockClient.EXPECT().ListSchema(mock.Anything).Return(schema1, nil)
			mockClient.EXPECT().ListTable(mock.Anything, "my_schema").Return(table1[:1], nil)
			mockClient.EXPECT().GetTableSchema(mock.Anything, table1[0]).Return("VIRTUAL_VIEW", schemaMapping[table1[0].Name()], nil)
			mockClient.EXPECT().GetMaskingPolicies(mock.Anything).Return(map[string][]string{}, nil)
		}, nil)

		require.NoError(t, err)
		require.Len(t, actual, 1)

		tableData := &v1beta2.Table{}
		require.NoError(t, proto.Unmarshal(actual[0].GetData().GetValue(), tableData))
		assert.Equal(t, "Dummy PDG", tableData.GetLabels()["pdg"])
	})

	t.Run("should not add pdg label when group mapping has no match for data_owner_team", func(t *testing.T) {
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
			mockClient.EXPECT().ListSchema(mock.Anything).Return(schema1, nil)
			mockClient.EXPECT().ListTable(mock.Anything, "my_schema").Return(table1[:1], nil)
			mockClient.EXPECT().GetTableSchema(mock.Anything, table1[0]).Return("VIRTUAL_VIEW", schemaMapping[table1[0].Name()], nil)
			mockClient.EXPECT().GetMaskingPolicies(mock.Anything).Return(map[string][]string{}, nil)
		}, nil)

		require.NoError(t, err)
		require.Len(t, actual, 1)

		tableData := &v1beta2.Table{}
		require.NoError(t, proto.Unmarshal(actual[0].GetData().GetValue(), tableData))
		assert.NotContains(t, tableData.GetLabels(), "pdg")
	})

	t.Run("should not set pdg label due to shield group response", func(t *testing.T) {
		tests := []struct {
			name        string
			shieldSetup func(t *testing.T) string
		}{
			{
				name: "when shield response has no groups key",
				shieldSetup: func(t *testing.T) string {
					s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("Content-Type", "application/json")
						_, _ = w.Write([]byte(`{"other_key": []}`))
					}))
					t.Cleanup(s.Close)
					return s.URL
				},
			},
			{
				name: "when group has no metadata field",
				shieldSetup: func(t *testing.T) string {
					return newGroupMappingServer(t, []map[string]interface{}{
						{"slug": "dummy_team"},
					}).URL
				},
			},
			{
				name: "when group metadata has no product-group-name",
				shieldSetup: func(t *testing.T) string {
					return newGroupMappingServer(t, []map[string]interface{}{
						{
							"slug": "dummy_team",
							"metadata": map[string]interface{}{
								"other-key": "some-value",
							},
						},
					}).URL
				},
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				actual, err := runTest(t, plugins.Config{
					URNScope: "test-maxcompute",
					RawConfig: map[string]interface{}{
						"project_name": projectID,
						"access_key": map[string]interface{}{
							"id":     "access_key_id",
							"secret": "access_key_secret",
						},
						"endpoint_project": "https://example.com/some-api",
						"shield_host":      tc.shieldSetup(t),
					},
				}, func(mockClient *mocks.MaxComputeClient) {
					mockClient.EXPECT().ListSchema(mock.Anything).Return(schema1, nil)
					mockClient.EXPECT().ListTable(mock.Anything, "my_schema").Return(table1[:1], nil)
					mockClient.EXPECT().GetTableSchema(mock.Anything, table1[0]).Return("VIRTUAL_VIEW", schemaMapping[table1[0].Name()], nil)
					mockClient.EXPECT().GetMaskingPolicies(mock.Anything).Return(map[string][]string{}, nil)
				}, nil)

				require.NoError(t, err)
				require.Len(t, actual, 1)

				tableData := &v1beta2.Table{}
				require.NoError(t, proto.Unmarshal(actual[0].GetData().GetValue(), tableData))
				assert.NotContains(t, tableData.GetLabels(), "pdg")
			})
		}
	})

	t.Run("should forward shield_header values as HTTP headers to Shield endpoint", func(t *testing.T) {
		var capturedAuthEmail string
		headerCapturingServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			capturedAuthEmail = r.Header.Get("X-Auth-Email")
			w.Header().Set("Content-Type", "application/json")
			body := map[string]interface{}{"groups": []interface{}{}}
			if err := json.NewEncoder(w).Encode(body); err != nil {
				http.Error(w, "encode error", http.StatusInternalServerError)
			}
		}))
		defer headerCapturingServer.Close()

		actual, err := runTest(t, plugins.Config{
			URNScope: "test-maxcompute",
			RawConfig: map[string]interface{}{
				"project_name": projectID,
				"access_key": map[string]interface{}{
					"id":     "access_key_id",
					"secret": "access_key_secret",
				},
				"endpoint_project": "https://example.com/some-api",
				"shield_host":      headerCapturingServer.URL,
				"shield_header": map[string]interface{}{
					"X-Auth-Email": "meteor-app@gojek.com",
				},
			},
		}, func(mockClient *mocks.MaxComputeClient) {
			mockClient.EXPECT().ListSchema(mock.Anything).Return(schema1, nil)
			mockClient.EXPECT().ListTable(mock.Anything, "my_schema").Return(table1[:1], nil)
			mockClient.EXPECT().GetTableSchema(mock.Anything, table1[0]).Return("VIRTUAL_VIEW", schemaMapping[table1[0].Name()], nil)
			mockClient.EXPECT().GetMaskingPolicies(mock.Anything).Return(map[string][]string{}, nil)
		}, nil)

		require.NoError(t, err)
		require.Len(t, actual, 1)
		assert.Equal(t, "meteor-app@gojek.com", capturedAuthEmail)
	})

	t.Run("should continue extraction with empty pdg mapping when listGroupMapping fails after retries", func(t *testing.T) {
		failingServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer failingServer.Close()

		actual, err := runTest(t, plugins.Config{
			URNScope: "test-maxcompute",
			RawConfig: map[string]interface{}{
				"project_name": projectID,
				"access_key": map[string]interface{}{
					"id":     "access_key_id",
					"secret": "access_key_secret",
				},
				"endpoint_project": "https://example.com/some-api",
				"shield_host":      failingServer.URL,
			},
		}, func(mockClient *mocks.MaxComputeClient) {
			mockClient.EXPECT().ListSchema(mock.Anything).Return(schema1, nil)
			mockClient.EXPECT().ListTable(mock.Anything, "my_schema").Return(table1[:1], nil)
			mockClient.EXPECT().GetTableSchema(mock.Anything, table1[0]).Return("VIRTUAL_VIEW", schemaMapping[table1[0].Name()], nil)
			mockClient.EXPECT().GetMaskingPolicies(mock.Anything).Return(map[string][]string{}, nil)
		}, nil)

		require.NoError(t, err)
		require.Len(t, actual, 1)

		tableData := &v1beta2.Table{}
		require.NoError(t, proto.Unmarshal(actual[0].GetData().GetValue(), tableData))
		assert.NotContains(t, tableData.GetLabels(), "pdg")
	})

	t.Run("should continue extraction when shield response has no groups key", func(t *testing.T) {
		noGroupsKeyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"other_key": []}`))
		}))
		defer noGroupsKeyServer.Close()

		actual, err := runTest(t, plugins.Config{
			URNScope: "test-maxcompute",
			RawConfig: map[string]interface{}{
				"project_name": projectID,
				"access_key": map[string]interface{}{
					"id":     "access_key_id",
					"secret": "access_key_secret",
				},
				"endpoint_project": "https://example.com/some-api",
				"shield_host":      noGroupsKeyServer.URL,
			},
		}, func(mockClient *mocks.MaxComputeClient) {
			mockClient.EXPECT().ListSchema(mock.Anything).Return(schema1, nil)
			mockClient.EXPECT().ListTable(mock.Anything, "my_schema").Return(table1[:1], nil)
			mockClient.EXPECT().GetTableSchema(mock.Anything, table1[0]).Return("VIRTUAL_VIEW", schemaMapping[table1[0].Name()], nil)
			mockClient.EXPECT().GetMaskingPolicies(mock.Anything).Return(map[string][]string{}, nil)
		}, nil)

		require.NoError(t, err)
		require.Len(t, actual, 1)

		tableData := &v1beta2.Table{}
		require.NoError(t, proto.Unmarshal(actual[0].GetData().GetValue(), tableData))
		assert.NotContains(t, tableData.GetLabels(), "pdg")
	})

	t.Run("should skip pdg enrichment when group has no metadata field", func(t *testing.T) {
		actual, err := runTest(t, plugins.Config{
			URNScope: "test-maxcompute",
			RawConfig: map[string]interface{}{
				"project_name": projectID,
				"access_key": map[string]interface{}{
					"id":     "access_key_id",
					"secret": "access_key_secret",
				},
				"endpoint_project": "https://example.com/some-api",
				"shield_host": newGroupMappingServer(t, []map[string]interface{}{
					{"slug": "dummy_team"},
				}).URL,
			},
		}, func(mockClient *mocks.MaxComputeClient) {
			mockClient.EXPECT().ListSchema(mock.Anything).Return(schema1, nil)
			mockClient.EXPECT().ListTable(mock.Anything, "my_schema").Return(table1[:1], nil)
			mockClient.EXPECT().GetTableSchema(mock.Anything, table1[0]).Return("VIRTUAL_VIEW", schemaMapping[table1[0].Name()], nil)
			mockClient.EXPECT().GetMaskingPolicies(mock.Anything).Return(map[string][]string{}, nil)
		}, nil)

		require.NoError(t, err)
		require.Len(t, actual, 1)

		tableData := &v1beta2.Table{}
		require.NoError(t, proto.Unmarshal(actual[0].GetData().GetValue(), tableData))
		assert.NotContains(t, tableData.GetLabels(), "pdg")
	})

	t.Run("should skip pdg enrichment when group metadata has no product-group-name", func(t *testing.T) {
		actual, err := runTest(t, plugins.Config{
			URNScope: "test-maxcompute",
			RawConfig: map[string]interface{}{
				"project_name": projectID,
				"access_key": map[string]interface{}{
					"id":     "access_key_id",
					"secret": "access_key_secret",
				},
				"endpoint_project": "https://example.com/some-api",
				"shield_host": newGroupMappingServer(t, []map[string]interface{}{
					{
						"slug": "dummy_team",
						"metadata": map[string]interface{}{
							"other-key": "some-value",
						},
					},
				}).URL,
			},
		}, func(mockClient *mocks.MaxComputeClient) {
			mockClient.EXPECT().ListSchema(mock.Anything).Return(schema1, nil)
			mockClient.EXPECT().ListTable(mock.Anything, "my_schema").Return(table1[:1], nil)
			mockClient.EXPECT().GetTableSchema(mock.Anything, table1[0]).Return("VIRTUAL_VIEW", schemaMapping[table1[0].Name()], nil)
			mockClient.EXPECT().GetMaskingPolicies(mock.Anything).Return(map[string][]string{}, nil)
		}, nil)

		require.NoError(t, err)
		require.Len(t, actual, 1)

		tableData := &v1beta2.Table{}
		require.NoError(t, proto.Unmarshal(actual[0].GetData().GetValue(), tableData))
		assert.NotContains(t, tableData.GetLabels(), "pdg")
	})

	t.Run("should skip schemas in the exclusion list", func(t *testing.T) {
		schemaA := odps.NewSchema(nil, projectID, "schema_a")
		schemaC := odps.NewSchema(nil, projectID, "schema_c")
		tableC := odps.NewTable(nil, projectID, "schema_c", "table_c")

		tableCSchemaBuilder := tableschema.NewSchemaBuilder()
		tableCSchemaBuilder.Name("table_c").Columns(c3)
		tableCSchema := tableCSchemaBuilder.Build()
		tableCSchema.TableName = "table_c"

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
					"schemas": []string{"schema_a"},
				},
			},
		}, func(mockClient *mocks.MaxComputeClient) {
			mockClient.EXPECT().ListSchema(mock.Anything).Return([]*odps.Schema{schemaA, schemaC}, nil)
			mockClient.EXPECT().ListTable(mock.Anything, "schema_c").Return([]*odps.Table{tableC}, nil)
			mockClient.EXPECT().GetTableSchema(mock.Anything, tableC).Return("MANAGED_TABLE", &tableCSchema, nil)
			mockClient.EXPECT().GetMaskingPolicies(mock.Anything).Return(map[string][]string{}, nil)
		}, nil)

		require.NoError(t, err)
		require.Len(t, actual, 1)
		assert.Equal(t, plugins.MaxComputeURN(projectID, "schema_c", "table_c"), actual[0].Urn)
	})

	t.Run("should only process schema matching schema_name filter", func(t *testing.T) {
		schemaA := odps.NewSchema(nil, projectID, "schema_a")
		schemaB := odps.NewSchema(nil, projectID, "schema_b")
		tableB := odps.NewTable(nil, projectID, "schema_b", "table_b")

		tableBSchemaBuilder := tableschema.NewSchemaBuilder()
		tableBSchemaBuilder.Name("table_b").Columns(c3)
		tableBSchema := tableBSchemaBuilder.Build()
		tableBSchema.TableName = "table_b"

		actual, err := runTest(t, plugins.Config{
			URNScope: "test-maxcompute",
			RawConfig: map[string]interface{}{
				"project_name": projectID,
				"access_key": map[string]interface{}{
					"id":     "access_key_id",
					"secret": "access_key_secret",
				},
				"endpoint_project": "https://example.com/some-api",
				"schema_name":      "schema_b",
			},
		}, func(mockClient *mocks.MaxComputeClient) {
			mockClient.EXPECT().ListSchema(mock.Anything).Return([]*odps.Schema{schemaA, schemaB}, nil)
			mockClient.EXPECT().ListTable(mock.Anything, "schema_b").Return([]*odps.Table{tableB}, nil)
			mockClient.EXPECT().GetTableSchema(mock.Anything, tableB).Return("MANAGED_TABLE", &tableBSchema, nil)
			mockClient.EXPECT().GetMaskingPolicies(mock.Anything).Return(map[string][]string{}, nil)
		}, nil)

		require.NoError(t, err)
		require.Len(t, actual, 1)
		assert.Equal(t, plugins.MaxComputeURN(projectID, "schema_b", "table_b"), actual[0].Urn)
	})

	t.Run("should skip schema and continue when ListTable fails", func(t *testing.T) {
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
			mockClient.EXPECT().ListSchema(mock.Anything).Return(schema1, nil)
			mockClient.EXPECT().ListTable(mock.Anything, "my_schema").Return(nil, fmt.Errorf("ListTable fails"))
		}, nil)

		assert.NoError(t, err)
		assert.Empty(t, actual)
	})

	t.Run("should continue extraction when GetMaskingPolicies fails", func(t *testing.T) {
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
			mockClient.EXPECT().ListSchema(mock.Anything).Return(schema1, nil)
			mockClient.EXPECT().ListTable(mock.Anything, "my_schema").Return(table1[:1], nil)
			mockClient.EXPECT().GetTableSchema(mock.Anything, table1[0]).Return("VIRTUAL_VIEW", schemaMapping[table1[0].Name()], nil)
			mockClient.EXPECT().GetMaskingPolicies(mock.Anything).Return(nil, fmt.Errorf("masking policy unavailable"))
		}, nil)

		require.NoError(t, err)
		require.Len(t, actual, 1)
	})

	t.Run("should not skip table with permanent lifecycle (-1) when min_table_lifecycle is set", func(t *testing.T) {
		permanentTable := odps.NewTable(nil, projectID, "my_schema", "permanent_table")

		permanentSchemaBuilder := tableschema.NewSchemaBuilder()
		permanentSchemaBuilder.Name("permanent_table").Columns(c3)
		permanentTableSchema := permanentSchemaBuilder.Build()
		permanentTableSchema.TableName = "permanent_table"
		permanentTableSchema.Lifecycle = -1

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
					"min_table_lifecycle": 5,
				},
			},
		}, func(mockClient *mocks.MaxComputeClient) {
			mockClient.EXPECT().ListSchema(mock.Anything).Return(schema1, nil)
			mockClient.EXPECT().ListTable(mock.Anything, "my_schema").Return([]*odps.Table{permanentTable}, nil)
			mockClient.EXPECT().GetTableSchema(mock.Anything, permanentTable).Return("MANAGED_TABLE", &permanentTableSchema, nil)
			mockClient.EXPECT().GetMaskingPolicies(mock.Anything).Return(map[string][]string{}, nil)
		}, nil)

		require.NoError(t, err)
		require.Len(t, actual, 1)
		assert.Equal(t, plugins.MaxComputeURN(projectID, "my_schema", "permanent_table"), actual[0].Urn)
	})

	t.Run("should mix preview values when mix_values is true", func(t *testing.T) {
		fixedRandomizer := func(seed int64) func(int64) int64 {
			return func(max int64) int64 {
				return max - 1
			}
		}

		actual, err := runTest(t, plugins.Config{
			URNScope: "test-maxcompute",
			RawConfig: map[string]interface{}{
				"project_name": projectID,
				"access_key": map[string]interface{}{
					"id":     "access_key_id",
					"secret": "access_key_secret",
				},
				"endpoint_project": "https://example.com/some-api",
				"max_preview_rows": 3,
				"mix_values":       true,
			},
		}, func(mockClient *mocks.MaxComputeClient) {
			mockClient.EXPECT().ListSchema(mock.Anything).Return(schema1, nil)
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
		}, fixedRandomizer)

		require.NoError(t, err)
		require.Len(t, actual, 1)

		tableData := &v1beta2.Table{}
		require.NoError(t, proto.Unmarshal(actual[0].GetData().GetValue(), tableData))
		require.NotNil(t, tableData.GetPreviewRows())

		rows := tableData.GetPreviewRows().GetValues()
		require.Len(t, rows, 2)

		firstRow := rows[0].GetListValue().GetValues()
		secondRow := rows[1].GetListValue().GetValues()

		assert.Equal(t, "2", firstRow[0].GetStringValue())
		assert.Equal(t, "user2@example.com", firstRow[1].GetStringValue())
		assert.Equal(t, "1", secondRow[0].GetStringValue())
		assert.Equal(t, "user1@example.com", secondRow[1].GetStringValue())
	})
}
