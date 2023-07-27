package bigtable

import (
	"context"
	"time"

	"cloud.google.com/go/bigtable"
	"github.com/goto/meteor/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type AdminClientMW struct {
	tableDuration  metric.Int64Histogram
	tablesDuration metric.Int64Histogram
	projectID      string
	instanceName   string
	next           AdminClient
}

type InstancesAdminClientMW struct {
	instancesDuration metric.Int64Histogram
	projectID         string
	next              InstanceAdminClient
}

func WithAdminClientMW(next AdminClient, projectID, instanceName string) (AdminClient, error) {
	meter := otel.Meter("")

	tablesDuration, err := meter.Int64Histogram("meteor.bigtable.client.tables.duration", metric.WithUnit("ms"))
	if err != nil {
		return nil, err
	}

	tableDuration, err := meter.Int64Histogram("meteor.bigtable.client.table.duration", metric.WithUnit("ms"))
	if err != nil {
		return nil, err
	}

	return &AdminClientMW{
		tableDuration:  tableDuration,
		tablesDuration: tablesDuration,
		projectID:      projectID,
		instanceName:   instanceName,
		next:           next,
	}, nil
}

func WithInstancesAdminClientMW(next InstanceAdminClient, projectID string) (InstanceAdminClient, error) {
	meter := otel.Meter("")

	instancesDuration, err := meter.Int64Histogram("meteor.bigtable.client.instances.duration", metric.WithUnit("ms"))
	if err != nil {
		return nil, err
	}

	return &InstancesAdminClientMW{
		instancesDuration: instancesDuration,
		projectID:         projectID,
		next:              next,
	}, nil
}

func (o *AdminClientMW) Tables(ctx context.Context) (res []string, err error) {
	defer func(start time.Time) {
		o.tablesDuration.Record(ctx,
			time.Since(start).Milliseconds(),
			metric.WithAttributes(
				attribute.String("bt.project_id", o.projectID),
				attribute.String("bt.error_code", utils.StatusCode(err).String()),
				attribute.String("bt.instance_name", o.instanceName),
			))
	}(time.Now())

	res, err = o.next.Tables(ctx)
	return res, err
}

func (o *AdminClientMW) TableInfo(ctx context.Context, table string) (res *bigtable.TableInfo, err error) {
	defer func(start time.Time) {
		o.tableDuration.Record(ctx,
			time.Since(start).Milliseconds(),
			metric.WithAttributes(
				attribute.String("bt.project_id", o.projectID),
				attribute.String("bt.error_code", utils.StatusCode(err).String()),
				attribute.String("bt.instance_name", o.instanceName),
				attribute.String("bt.table_name", table),
			))
	}(time.Now())
	res, err = o.next.TableInfo(ctx, table)
	return res, err
}

func (o *InstancesAdminClientMW) Instances(ctx context.Context) (res []*bigtable.InstanceInfo, err error) {
	defer func(start time.Time) {
		o.instancesDuration.Record(ctx,
			time.Since(start).Milliseconds(),
			metric.WithAttributes(
				attribute.String("bt.project_id", o.projectID),
				attribute.String("bt.error_code", utils.StatusCode(err).String()),
			))
	}(time.Now())

	res, err = o.next.Instances(ctx)
	return res, err
}
