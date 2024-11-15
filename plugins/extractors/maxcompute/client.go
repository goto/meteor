package maxcompute

import (
	"context"

	"github.com/aliyun/aliyun-odps-go-sdk/odps"
	"github.com/aliyun/aliyun-odps-go-sdk/odps/account"
)

type MCClient struct {
	client  *odps.Odps
	project *odps.Project
}

func NewMaxComputeClient(config Config) *MCClient {
	aliAccount := account.NewAliyunAccount(config.AccessKey.ID, config.AccessKey.Secret)
	client := odps.NewOdps(aliAccount, config.EndpointProject)
	client.SetDefaultProjectName(config.ProjectName)

	project := client.Project(config.ProjectName)

	return &MCClient{
		client:  client,
		project: project,
	}
}

func (c *MCClient) ListSchema(context.Context) (schemas []*odps.Schema, err error) {
	err = c.project.Schemas().List(func(schema *odps.Schema, err2 error) {
		if err2 != nil {
			err = err2
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
				return
			}
			tables = append(tables, table)
		},
	)
	return tables, err
}

func (*MCClient) GetTable(_ context.Context, table *odps.Table) (*odps.Table, error) {
	err := table.Load()
	if err != nil {
		isView := table.Schema().IsVirtualView || table.Schema().IsMaterializedView
		isLoaded := table.IsLoaded()
		if !isView || (isView && !isLoaded) {
			return nil, err
		}
	}
	return table, nil
}
