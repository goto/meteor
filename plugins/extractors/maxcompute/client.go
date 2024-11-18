package maxcompute

import (
	"context"

	"github.com/aliyun/aliyun-odps-go-sdk/odps"
	"github.com/aliyun/aliyun-odps-go-sdk/odps/account"
	"github.com/aliyun/aliyun-odps-go-sdk/odps/tableschema"
	"github.com/goto/meteor/plugins"
	"github.com/goto/salt/log"
)

type MCClient struct {
	client  *odps.Odps
	project *odps.Project
	log     log.Logger
}

func NewMaxComputeClient(config Config) *MCClient {
	aliAccount := account.NewAliyunAccount(config.AccessKey.ID, config.AccessKey.Secret)
	client := odps.NewOdps(aliAccount, config.EndpointProject)
	client.SetDefaultProjectName(config.ProjectName)

	project := client.Project(config.ProjectName)

	return &MCClient{
		client:  client,
		project: project,
		log:     plugins.GetLog(),
	}
}

func (c *MCClient) ListSchema(context.Context) (schemas []*odps.Schema, err error) {
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

func (c *MCClient) ListTable(_ context.Context, schemaName string) (tables []*odps.Table, err error) {
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

func (*MCClient) GetTableSchema(_ context.Context, table *odps.Table) (string, *tableschema.TableSchema, error) {
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
