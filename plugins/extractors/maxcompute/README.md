# maxcompute
## Usage
## TODO: Fix the description

```yaml
source:
    name: maxcompute
    config:
        project_name: goto_test
        endpoint_project: http://goto_test-maxcompute.com
        access_key_json: {
            id: xyz
            secret_base64: __base64__
        }
        schema_name: DEFAULT
        max_preview_rows: 3
        exclude:
            schemas:
                - schema_a
                - schema_b
            tables:
                - schema_c.table_a
        max_page_size: 100
        concurrency: 10
        mix_values: false
        include_column_profile: true
        build_view_lineage: true
        collect_table_usage: false
        usage_period_in_day: 7
        usage_project_names:
            - maxcompute-project-name
            - other-maxcompute-project-name
```

## Inputs

| Key | Value | Example | Description |    |
| :-- | :---- | :------ | :---------- | :-- |
| `project_name` | `string` | `goto_test` | MaxCompute Project Name | *required* |
| `endpoint_project` | `string` | `http://goto_test-maxcompute.com` | Endpoint Project URL | *required* |
| `access_key_json` | `string` | `{"id": "xyz", "secret_base64": "__base64__"}` | Access Key in JSON string | *required* |
| `schema_name` | `string` | `DEFAULT` | Default schema name | *optional* |
| `max_preview_rows` | `int` | `3` | max number of preview rows to fetch, `0` will skip preview fetching. Default to `3`. | *optional* |
| `exclude.schemas` | `[]string` | `["schema_a", "schema_b"]` | List of schemas to exclude | *optional* |
| `exclude.tables` | `[]string` | `["schema_c.table_a"]` | List of tables to exclude | *optional* |
| `max_page_size` | `int` | `100` | max page size hint used for fetching datasets/tables/rows from MaxCompute | *optional* |
| `concurrency` | `int` | `10` | Number of concurrent requests to MaxCompute | *optional* |
| `mix_values` | `bool` | `false` | true if you want to mix the column values with the preview rows. Default to `false`. | *optional* |
| `include_column_profile` | `bool` | `true` | true if you want to profile the column value such min, max, med, avg, top, and freq | *optional* |
| `build_view_lineage` | `bool` | `true` | true if you want to build view lineage | *optional* |
| `collect_table_usage` | `boolean` | `false` | toggle feature to collect table usage, `true` will enable collecting table usage. Default to `false`. | *optional* |
| `usage_period_in_day` | `int` | `7` | collecting log from `(now - usage_period_in_day)` until `now`. only matter if `collect_table_usage` is true. Default to `7`. | *optional* |
| `usage_project_names` | `[]string` | `["maxcompute-project-name", "other-maxcompute-project-name"]` | collecting log from defined MaxCompute Project Names. | *optional* |

### *Notes*

- Leaving `access_key_json` blank will default to [MaxCompute's default authentication][maxcompute-default-auth].
- Setting `max_preview_rows` to `0` will skip preview fetching.

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

