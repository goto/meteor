package client

import (
	"context"

	"github.com/aliyun/aliyun-odps-go-sdk/odps"
	"github.com/aliyun/aliyun-odps-go-sdk/odps/account"
	"github.com/aliyun/aliyun-odps-go-sdk/odps/tableschema"
	"github.com/aliyun/aliyun-odps-go-sdk/odps/tunnel"
	"github.com/goto/meteor/plugins/extractors/maxcompute/config"
	"google.golang.org/protobuf/types/known/structpb"
)

type Client struct {
	client  *odps.Odps
	project *odps.Project
	tunnel  *tunnel.Tunnel

	isSchemaEnabled bool
}

func New(conf config.Config) (*Client, error) {
	aliAccount := account.NewAliyunAccount(conf.AccessKey.ID, conf.AccessKey.Secret)
	client := odps.NewOdps(aliAccount, conf.EndpointProject)
	client.SetDefaultProjectName(conf.ProjectName)

	project := client.Project(conf.ProjectName)
	tunnelInstance, err := tunnel.NewTunnelFromProject(project)
	if err != nil {
		return nil, err
	}

	properties, err := project.GetAllProperties()
	if err != nil {
		return nil, err
	}

	const schemaProperty = "odps.schema.model.enabled"
	isSchemaEnabled := properties.Get(schemaProperty) == "true"

	return &Client{
		client:          client,
		project:         project,
		tunnel:          tunnelInstance,
		isSchemaEnabled: isSchemaEnabled,
	}, nil
}

func (c *Client) ListSchema(_ context.Context) (schemas []*odps.Schema, err error) {
	if !c.isSchemaEnabled {
		schema := odps.NewSchema(nil, "", "default")
		return []*odps.Schema{schema}, nil
	}

	err = c.project.Schemas().List(func(schema *odps.Schema, err2 error) {
		if err2 != nil {
			err = err2
			return
		}
		schemas = append(schemas, schema)
	})

	return schemas, err
}

func (c *Client) ListTable(_ context.Context, schemaName string) (tables []*odps.Table, err error) {
	t := odps.NewTables(c.client, c.project.Name(), schemaName)
	t.List(
		func(table *odps.Table, err2 error) {
			if err2 != nil {
				err = err2
				return
			}
			tables = append(tables, table)
		},
	)
	return tables, err
}

func (*Client) GetTableSchema(_ context.Context, table *odps.Table) (string, *tableschema.TableSchema, error) {
	err := table.Load()
	tableSchema := table.Schema()
	if err != nil {
		isView := tableSchema.IsVirtualView || tableSchema.IsMaterializedView
		isLoaded := table.IsLoaded()
		if !isView || (isView && !isLoaded) {
			return "", nil, err
		}
	}
	return table.Type().String(), &tableSchema, nil
}

func (c *Client) GetTablePreview(_ context.Context, partitionValue string, table *odps.Table, maxRows int) (
	previewFields []string, previewRows *structpb.ListValue, err error,
) {
	if table.Type().String() == config.TableTypeView {
		return nil, nil, nil
	}

	records, err := c.tunnel.Preview(table, partitionValue, int64(maxRows))
	if err != nil {
		return nil, nil, err
	}

	columnNames := make([]string, len(table.Schema().Columns))
	for i, column := range table.Schema().Columns {
		columnNames[i] = column.Name
	}

	var protoRows []*structpb.Value
	for _, record := range records {
		var rowValues []*structpb.Value
		for _, value := range record {
			rowValues = append(rowValues, structpb.NewStringValue(value.String()))
		}
		protoRows = append(protoRows, structpb.NewListValue(&structpb.ListValue{Values: rowValues}))
	}

	protoList := &structpb.ListValue{Values: protoRows}

	return columnNames, protoList, nil
}

func (c *Client) GetPolicyTagsAndMaskingPolicy(table *odps.Table) (string, []string, error) {
	var policyTags []string
	var maskingPolicy string
	colPolicyTags, err := table.ColumnMaskInfos()
	if err != nil {
		return maskingPolicy, nil, err
	}
	for _, policyTag := range colPolicyTags {
		maskingPolicy = policyTag.Name
		policyTags = append(policyTags, policyTag.PolicyNameList...)
	}
	return maskingPolicy, policyTags, nil
}
