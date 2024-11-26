package client

import (
	"context"

	"github.com/aliyun/aliyun-odps-go-sdk/odps"
	"github.com/aliyun/aliyun-odps-go-sdk/odps/account"
	"github.com/aliyun/aliyun-odps-go-sdk/odps/tableschema"
	"github.com/aliyun/aliyun-odps-go-sdk/odps/tunnel"
	"github.com/goto/meteor/plugins"
	"github.com/goto/meteor/plugins/extractors/maxcompute/config"
	"github.com/goto/salt/log"
	"google.golang.org/protobuf/types/known/structpb"
)

type Client struct {
	client  *odps.Odps
	project *odps.Project
	tunnel  *tunnel.Tunnel
	log     log.Logger
}

func New(conf config.Config) *Client {
	aliAccount := account.NewAliyunAccount(conf.AccessKey.ID, conf.AccessKey.Secret)
	client := odps.NewOdps(aliAccount, conf.EndpointProject)
	client.SetDefaultProjectName(conf.ProjectName)

	project := client.Project(conf.ProjectName)
	tunnelInstance, err := tunnel.NewTunnelFromProject(project)
	if err != nil {
		plugins.GetLog().Error("failed to create tunnel", "with error:", err)
	}

	return &Client{
		client:  client,
		project: project,
		log:     plugins.GetLog(),
		tunnel:  tunnelInstance,
	}
}

func (c *Client) ListSchema(context.Context) (schemas []*odps.Schema, err error) {
	err = c.project.Schemas().List(func(schema *odps.Schema, err2 error) {
		if err2 != nil {
			err = err2
			c.log.Error("failed to process schema", "with error:", err)
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
				c.log.Error("failed to process table", "with error:", err)
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

func (c *Client) GetTablePreview(ctx context.Context, partitionValue string, table *odps.Table, maxRows int) ([]string, *structpb.ListValue, error) {
	records, err := c.tunnel.Preview(table, partitionValue, int64(maxRows))
	if err != nil {
		c.log.Error("failed to preview table", "table", table.Name(), "error", err)
		return nil, nil, err
	}

	columnNames := make([]string, len(records))
	for i, col := range records {
		columnNames[i] = col.String()
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
