# maxcompute
## Usage
## TODO: Fix the description

```yaml
source:
    name: maxcompute
    config:
        project_name: goto_test
        endpoint_project: http://goto_test-maxcompute.com
        access_key:
            id: access_key_id
            secret: access_key_secret
        schema_name: DEFAULT
        exclude:
            schemas:
                - schema_a
                - schema_b
            tables:
                - schema_c.table_a
        concurrency: 10
        build_view_lineage: true
```

## Inputs

| Key | Value | Example | Description |    |
| :-- | :---- | :------ | :---------- | :-- |
| `project_name` | `string` | `goto_test` | MaxCompute Project Name | *required* |
| `endpoint_project` | `string` | `http://goto_test-maxcompute.com` | Endpoint Project URL | *required* |
| `access_key.id` | `string` | `access_key_id` | Access Key ID | *required* |
| `access_key.secret` | `string` | `access_key_secret` | Access Key Secret | *required* |
| `schema_name` | `string` | `DEFAULT` | Default schema name | *optional* |
| `exclude.schemas` | `[]string` | `["schema_a", "schema_b"]` | List of schemas to exclude | *optional* |
| `exclude.tables` | `[]string` | `["schema_c.table_a"]` | List of tables to exclude | *optional* |
| `concurrency` | `int` | `10` | Number of concurrent requests to MaxCompute | *optional* |
| `build_view_lineage` | `bool` | `true` | true if you want to build view lineage | *optional* |

### *Notes*

- Leaving `access_key` blank will default to [MaxCompute's default authentication][maxcompute-default-auth].

## Outputs

| Field                          | Sample Value                                                                                                                                                                                                                                                                                                                                                                                | Description                                             |
|:-------------------------------|:--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|:--------------------------------------------------------|
| `resource.urn`                 | `project_name.schema_name.table_name`                                                                                                                                                                                                                                                                                                                                                        |                                                         |
| `resource.name`                | `table_name`                                                                                                                                                                                                                                                                                                                                                                                |                                                         |
| `resource.service`             | `maxcompute`                                                                                                                                                                                                                                                                                                                                                                                  |                                                         |
| `description`                  | `table description`                                                                                                                                                                                                                                                                                                                                                                         |                                                         |
| `profile.total_rows`           | `2100`                                                                                                                                                                                                                                                                                                                                                                                      |                                                         |
| `profile.usage_count`          | `15`                                                                                                                                                                                                                                                                                                                                                                                        |                                                         |
| `profile.joins`                | [][Join](#Join)                                                                                                                                                                                                                                                                                                                                                                             |                                                         |
| `profile.filters`              | [`"WHERE t.param_3 = 'the_param' AND t.column_1 = \"xxxxxx-xxxx-xxxx-xxxx-xxxxxxxxx\""`,`"WHERE event_timestamp >= TIMESTAMP(\"2021-10-29\", \"UTC\") AND event_timestamp < TIMESTAMP(\"2021-11-22T02:01:06Z\")"`]                                                                                                                                                                          |                                                         |
| `schema`                       | [][Column](#column)                                                                                                                                                                                                                                                                                                                                                                         |                                                         |
| `properties.partition_data`    | `"partition_data": {"partition_field": "data_date", "require_partition_filter": false, "time_partition": {"partition_by": "DAY","partition_expire": 0 } }`                                                                                                                                                                                                                                  | partition related data for time and range partitioning. |
| `properties.clustering_fields` | `['created_at', 'updated_at']`                                                                                                                                                                                                                                                                                                                                                              | list of fields on which the table is clustered          |
| `properties.partition_field`   | `created_at`                                                                                                                                                                                                                                                                                                                                                                                | returns the field on which table is time partitioned    |

### Partition Data

| Field                                     | Sample Value | Description                                                                                                                                                                                                                                                             |
|:------------------------------------------|:-------------|:------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `partition_field`                         | `created_at` | field on which the table is partitioned either by TimePartitioning or RangePartitioning. In case field is empty for TimePartitioning _PARTITIONTIME is returned instead of empty.                                                                                       |
| `require_partition_filter`                | `true`       | boolean value which denotes if every query on the MaxCompute table must include at least one predicate that only references the partitioning column                                                                                                                       |
| `time_partition.partition_by`             | `HOUR`       | returns partition type HOUR/DAY/MONTH/YEAR                                                                                                                                                                                                                              |
| `time_partition.partition_expire_seconds` | `0`          | time in which data will expire from this partition. If 0 it will not expire.                                                                                                                                                                                            |
| `range_partition.interval`                | `10`         | width of a interval range                                                                                                                                                                                                                                               |
| `range_partition.start`                   | `0`          | start value for partition inclusive of this value                                                                                                                                                                                                                       |
| `range_partition.end`                     | `100`        | end value for partition exclusive of this value                                                                                                                                                                                                                         |


### Column

| Field         | Sample Value                           |
|:--------------|:---------------------------------------|
| `name`        | `total_price`                          |
| `description` | `item's total price`                   |
| `data_type`   | `decimal`                              |
| `is_nullable` | `true`                                 |
| `length`      | `12,2`                                 |
| `profile`     | `{"min":...,"max": ...,"unique": ...}` |

### Join

| Field        | Sample Value                                                                                                                                            |
|:-------------|:--------------------------------------------------------------------------------------------------------------------------------------------------------|
| `urn`        | `project_name.schema_name.table_name`                                                                                                                    |
| `count`      | `3`                                                                                                                                                     |
| `conditions` | [`"ON target.column_1 = source.column_1 and target.param_name = source.param_name"`,`"ON DATE(target.event_timestamp) = DATE(source.event_timestamp)"`] |

## Contributing

Refer to the [contribution guidelines](../../../docs/docs/contribute/guide.md#adding-a-new-extractor) for information on 
contributing to this module.

[maxcompute-default-auth]: https://www.alibabacloud.com/help/doc-detail/27800.htm

