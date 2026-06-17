package maxcompute

import (
	"context"
	_ "embed" // used to print the embedded assets
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync/atomic"
	"time"

	"github.com/aliyun/aliyun-odps-go-sdk/odps"
	"github.com/aliyun/aliyun-odps-go-sdk/odps/common"
	"github.com/aliyun/aliyun-odps-go-sdk/odps/datatype"
	"github.com/aliyun/aliyun-odps-go-sdk/odps/tableschema"
	"github.com/goto/meteor/metrics/otelhttpclient"
	"github.com/goto/meteor/models"
	v1beta2 "github.com/goto/meteor/models/gotocompany/assets/v1beta2"
	"github.com/goto/meteor/plugins"
	"github.com/goto/meteor/plugins/extractors/maxcompute/client"
	"github.com/goto/meteor/plugins/extractors/maxcompute/config"
	"github.com/goto/meteor/plugins/internal/upstream"
	"github.com/goto/meteor/plugins/internal/urlbuilder"
	"github.com/goto/meteor/registry"
	"github.com/goto/meteor/utils"
	"github.com/goto/salt/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	maxcomputeService             = "maxcompute"
	attributesDataSQL             = "sql"
	attributesDataMaskingPolicy   = "masking_policy"
	attributesDataProjectName     = "project_name"
	attributesDataSchema          = "schema"
	attributesDataType            = "type"
	attributesDataResourceURL     = "resource_url"
	attributesDataPartitionFields = "partition_fields"
	attributesDataLabel           = "label"
	attributesDataLifecycle       = "lifecycle"
	attributesDataDDLStatement    = "ddl_statement"

	httpTimeout                = 30 * time.Second
	listGroupMappingRoute      = "/admin/v1beta1/groups"
	listGroupMappingMaxRetry   = 3
	listGroupMappingRetryDelay = 200 * time.Millisecond
)

type Extractor struct {
	plugins.BaseExtractor
	logger log.Logger
	config config.Config
	randFn randFn

	client     Client
	newClient  NewClientFunc
	urlb       urlbuilder.Source
	httpClient *http.Client
}

type randFn func(rndSeed int64) func(int64) int64

type NewClientFunc func(ctx context.Context, logger log.Logger, conf config.Config) (Client, error)

type extractionRun struct {
	groupMapping  map[string]string
	emit          plugins.Emit
	tablesSkipped *atomic.Int64
}

//go:embed README.md
var summary string

var sampleConfig = `
project_name: goto_test
endpoint_project: https://service.ap-southeast-5.maxcompute.aliyun.com/api
access_key:
    id: access_key_id
    secret: access_key_secret
schema_name: default
exclude:
	schemas:
    	- schema_a
    	- schema_b
    tables:
    	- schema_c.table_a
concurrency: 10
max_preview_rows: 3
mix_values: false
build_view_lineage: true,
`

var info = plugins.Info{
	Description:  "MaxCompute metadata and metrics",
	SampleConfig: sampleConfig,
	Tags:         []string{"alicloud", "table"},
	Summary:      summary,
}

//go:generate mockery --name=Client -r --case underscore --with-expecter --structname MaxComputeClient --filename maxcompute_client_mock.go --output=./mocks
type Client interface {
	ListSchema(ctx context.Context) ([]*odps.Schema, error)
	ListTable(ctx context.Context, schemaName string) ([]*odps.Table, error)
	GetTableSchema(ctx context.Context, table *odps.Table) (string, *tableschema.TableSchema, error)
	GetTablePreview(ctx context.Context, partitionValue string, table *odps.Table, maxRows int) ([]string, *structpb.ListValue, error)
	GetMaskingPolicies(table *odps.Table) (map[client.Column][]client.Policy, error)
}

func New(logger log.Logger, clientFunc NewClientFunc, randFn randFn) *Extractor {
	e := &Extractor{
		logger:    logger,
		newClient: clientFunc,
		randFn:    randFn,
	}
	e.BaseExtractor = plugins.NewBaseExtractor(info, &e.config)
	e.ScopeNotRequired = true

	return e
}

func (e *Extractor) Init(ctx context.Context, conf plugins.Config) error {
	if err := e.BaseExtractor.Init(ctx, conf); err != nil {
		return err
	}

	if e.config.ProjectName == "" {
		return fmt.Errorf("project_name is required")
	}
	if e.config.AccessKey.ID == "" || e.config.AccessKey.Secret == "" {
		return fmt.Errorf("access_key is required")
	}
	if e.config.EndpointProject == "" {
		return fmt.Errorf("endpoint_project is required")
	}
	if e.config.Concurrency == 0 {
		e.config.Concurrency = 1
	}

	var err error
	e.client, err = e.newClient(ctx, e.logger, e.config)
	if err != nil {
		return err
	}

	e.urlb, err = urlbuilder.NewSource(e.config.ShieldHost)
	if err != nil {
		return err
	}

	httpClient := &http.Client{Timeout: httpTimeout}
	httpClient.Transport = otelhttpclient.NewHTTPTransport(httpClient.Transport)
	e.httpClient = httpClient

	return nil
}

func (e *Extractor) Extract(ctx context.Context, emit plugins.Emit) error {
	groupMapping := e.listGroupMapping(ctx)

	schemas, err := e.client.ListSchema(ctx)
	if err != nil && len(schemas) == 0 {
		return err
	}

	var schemasSkipped int
	var tablesSkipped atomic.Int64
	run := &extractionRun{
		groupMapping:  groupMapping,
		emit:          emit,
		tablesSkipped: &tablesSkipped,
	}

	for _, schema := range schemas {
		if e.config.SchemaName != "" && schema.Name() != e.config.SchemaName {
			continue
		}
		if contains(e.config.Exclude.Schemas, schema.Name()) {
			e.logger.Info("skipping schema as it is in the exclude list", "schema", schema.Name())
			continue
		}

		if err := e.fetchTablesFromSchema(ctx, schema, run); err != nil {
			e.logger.Warn("skipping schema due to error", "schema", schema.Name(), "error", err)
			schemasSkipped++
			continue
		}
	}

	if schemasSkipped > 0 || tablesSkipped.Load() > 0 {
		e.logger.Warn("extraction completed with partial failures", "project", e.config.ProjectName,
			"schemas_skipped", schemasSkipped, "tables_skipped", tablesSkipped.Load())
	}

	return nil
}

func (e *Extractor) listTablesCopy(ctx context.Context, schemaName string) ([]*odps.Table, error) {
	rawTables, err := e.client.ListTable(ctx, schemaName)
	if err != nil && len(rawTables) == 0 {
		return nil, err
	}

	tables := make([]*odps.Table, len(rawTables))
	copy(tables, rawTables)
	return tables, nil
}

func (e *Extractor) fetchTablesFromSchema(ctx context.Context, schema *odps.Schema, run *extractionRun) error {
	tables, err := e.listTablesCopy(ctx, schema.Name())
	if err != nil {
		return err
	}

	eg := &errgroup.Group{}
	eg.SetLimit(e.config.Concurrency)

	for i, table := range tables {
		tableName := fmt.Sprintf("%s.%s", schema.Name(), table.Name())
		if contains(e.config.Exclude.Tables, tableName) {
			e.logger.Info("skipping table as it is in the exclude list", "table", tableName)
			tables[i] = nil
			continue
		}

		tbl := table
		tables[i] = nil
		eg.Go(func() error {
			return e.processTable(ctx, schema, tbl, run)
		})
	}

	if err := eg.Wait(); err != nil {
		return fmt.Errorf("process tables: %w", err)
	}

	return nil
}

func (e *Extractor) processTable(ctx context.Context, schema *odps.Schema, table *odps.Table, run *extractionRun) error {
	tableName := fmt.Sprintf("%s.%s", schema.Name(), table.Name())
	tableType, tableSchema, err := e.client.GetTableSchema(ctx, table)
	if err != nil {
		e.logger.Error("failed to get table schema", "table", tableName, "error", err)
		run.tablesSkipped.Add(1)
		return nil
	}

	// If the lifecycle is less than the minimum lifecycle (days), skip the table
	if e.config.Exclude.MinTableLifecycle > 1 {
		lifecyclePermanent := tableSchema.Lifecycle == -1
		lifecycleNotConfigured := tableSchema.Lifecycle == 0
		if !lifecyclePermanent && !lifecycleNotConfigured && tableSchema.Lifecycle < e.config.Exclude.MinTableLifecycle {
			e.logger.Info("skipping table due to lifecycle less than minimum configured lifecycle",
				"table", tableName, "lifecycle", tableSchema.Lifecycle)
			run.tablesSkipped.Add(1)
			return nil
		}
	}

	asset, err := e.buildAsset(ctx, schema, table, tableType, tableSchema, run.groupMapping)
	if err != nil {
		e.logger.Error("failed to build asset", "table", tableName, "error", err)
		run.tablesSkipped.Add(1)
		return nil
	}

	run.emit(models.NewRecord(asset))
	return nil
}

func (e *Extractor) buildAsset(ctx context.Context, schema *odps.Schema,
	table *odps.Table, tableType string, tableSchema *tableschema.TableSchema,
	groupMapping map[string]string,
) (*v1beta2.Asset, error) {
	defaultSchema := "default"
	schemaName := schema.Name()
	if schemaName == "" {
		schemaName = defaultSchema
	}

	tableURN := plugins.MaxComputeURN(e.config.ProjectName, schemaName, tableSchema.TableName)

	previewFields, previewRows, err := e.buildPreview(ctx, table, tableSchema)
	if err != nil {
		e.logger.Warn("error building preview", "err", err, "table", tableSchema.TableName)
	}

	asset := &v1beta2.Asset{
		Urn:         tableURN,
		Name:        tableSchema.TableName,
		Type:        "table",
		Description: tableSchema.Comment,
		CreateTime:  timestamppb.New(time.Time(tableSchema.CreateTime)),
		UpdateTime:  timestamppb.New(time.Time(tableSchema.LastModifiedTime)),
		Service:     maxcomputeService,
	}

	tableAttributesData := e.buildTableAttributesData(schemaName, tableType, tableSchema)

	if tableType == config.TableTypeView {
		query := tableSchema.ViewText
		tableAttributesData[attributesDataSQL] = query
		if e.config.BuildViewLineage {
			upstreamResources := getUpstreamResources(query)
			asset.Lineage = &v1beta2.Lineage{
				Upstreams: upstreamResources,
			}
		}
	}

	maskingPolicy, err := e.client.GetMaskingPolicies(table)
	if err != nil {
		e.logger.Warn("error getting masking policy", "error", err)
	}

	var columns []*v1beta2.Column
	for i, col := range tableSchema.Columns {
		colAttrs, err := utils.TryParseMapToProto(buildColumnAttributesData(&tableSchema.Columns[i]))
		if err != nil {
			e.logger.Warn("error building column attributes, using empty attributes", "column", col.Name, "error", err)
			colAttrs = &structpb.Struct{}
		}
		columnData := &v1beta2.Column{
			Name:        col.Name,
			DataType:    dataTypeToString(col.Type),
			Description: col.Comment,
			IsNullable:  !col.NotNull,
			Attributes:  colAttrs,
			Columns:     buildColumns(col.Type),
		}

		if policies, found := maskingPolicy[col.Name]; found {
			policyValues := make([]*structpb.Value, 0, len(policies))
			for _, policy := range policies {
				policyValues = append(policyValues, structpb.NewStringValue(policy))
			}
			columnData.Attributes.Fields[attributesDataMaskingPolicy] = &structpb.Value{
				Kind: &structpb.Value_ListValue{
					ListValue: &structpb.ListValue{Values: policyValues},
				},
			}
		}

		columns = append(columns, columnData)
	}

	tableProfile := &v1beta2.TableProfile{}
	if tableSchema.RecordNum >= 0 {
		tableProfile.TotalRows = int64(tableSchema.RecordNum)
	}

	tableAttrs, err := utils.TryParseMapToProto(tableAttributesData)
	if err != nil {
		e.logger.Warn("error building table attributes, using empty attributes", "table", tableSchema.TableName, "error", err)
		tableAttrs = &structpb.Struct{}
	}

	tableData := &v1beta2.Table{
		Profile:    tableProfile,
		Attributes: tableAttrs,
		Columns:    columns,
		CreateTime: timestamppb.New(time.Time(tableSchema.CreateTime)),
		UpdateTime: timestamppb.New(time.Time(tableSchema.LastModifiedTime)),
	}

	var decodedComment map[string]string
	comment := strings.NewReplacer("\n", `\n`, "\r", `\r`, "\t", `\t`).Replace(tableSchema.Comment)
	if err := json.Unmarshal([]byte(comment), &decodedComment); err != nil {
		e.logger.Warn("error decoding comment", "error", err)
	}

	if decodedComment != nil {
		if desc, ok := decodedComment["description"]; ok {
			asset.Description = desc
		}
		if shieldTeam, ok := decodedComment["data_owner_team"]; ok {
			if pdgName, ok := groupMapping[shieldTeam]; ok {
				decodedComment["pdg"] = pdgName
			}
		}
		tableData.Labels = decodedComment
	}

	maxPreviewRows := e.config.MaxPreviewRows
	if maxPreviewRows > 0 {
		tableData.PreviewFields = previewFields
		if e.config.IncludePreviewRows {
			tableData.PreviewRows = previewRows
		}
	}

	ddl, err := getDDLStatement(tableSchema, e.config.ProjectName, schemaName)
	if err != nil {
		e.logger.Warn("error generating DDL", "error", err, "table", tableSchema.TableName)
	} else {
		tableData.DdlStatement = ddl
	}

	tbl, err := anypb.New(tableData)
	if err != nil {
		e.logger.Warn("error creating Any struct", "error", err)
	}
	asset.Data = tbl

	return asset, nil
}

func (e *Extractor) listGroupMapping(ctx context.Context) map[string]string {
	if e.config.ShieldHost == "" {
		return map[string]string{}
	}

	var lastErr error
	for attempt := 1; attempt <= listGroupMappingMaxRetry; attempt++ {
		result, err := e.doListGroupMapping(ctx)
		if err == nil {
			return result
		}

		lastErr = err
		e.logger.Warn("group mapping fetch attempt failed, retrying",
			"attempt", attempt, "max_retries", listGroupMappingMaxRetry, "error", err)
		if attempt < listGroupMappingMaxRetry {
			select {
			case <-ctx.Done():
				e.logger.Warn("context canceled while retrying group mapping fetch", "error", ctx.Err())
				return map[string]string{}
			case <-time.After(listGroupMappingRetryDelay):
			}
		}
	}

	e.logger.Warn("failed to fetch group mapping after retries, continuing with empty PDG mapping", "error", lastErr)
	return map[string]string{}
}

func (e *Extractor) doListGroupMapping(ctx context.Context) (map[string]string, error) {
	targetURL := e.urlb.New().Path(listGroupMappingRoute).URL()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL.String(), nil)
	if err != nil {
		return nil, err
	}

	for key, value := range e.config.ShieldHeader {
		req.Header.Add(key, value)
	}
	req = otelhttpclient.AnnotateRequest(req, listGroupMappingRoute)

	res, err := e.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer plugins.DrainBody(res)

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response with status: %d", res.StatusCode)
	}

	var groupMapping map[string][]interface{}
	if err := json.NewDecoder(res.Body).Decode(&groupMapping); err != nil {
		return nil, err
	}

	groupDetails, found := groupMapping["groups"]
	if !found {
		return map[string]string{}, nil
	}

	groupResult := make(map[string]string)
	for _, group := range groupDetails {
		slug, pdg, found := extractGroupPDG(group)
		if found {
			groupResult[slug] = pdg
		}
	}

	return groupResult, nil
}

func extractGroupPDG(group interface{}) (slug, pdg string, found bool) {
	groupDetail, ok := group.(map[string]interface{})
	if !ok {
		return "", "", false
	}
	slug, ok = groupDetail["slug"].(string)
	if !ok {
		return "", "", false
	}
	metadata, found := groupDetail["metadata"]
	if !found {
		return "", "", false
	}
	metadataMap, ok := metadata.(map[string]interface{})
	if !ok {
		return "", "", false
	}
	pdgVal, found := metadataMap["product-group-name"]
	if !found {
		return "", "", false
	}
	pdg, ok = pdgVal.(string)
	return slug, pdg, ok
}

func getUpstreamResources(query string) []*v1beta2.Resource {
	upstreamDependencies := upstream.ParseTopLevelUpstreamsFromQuery(query)
	uniqueUpstreamDependencies := upstream.UniqueFilterResources(upstreamDependencies)
	var upstreams []*v1beta2.Resource
	for _, dependency := range uniqueUpstreamDependencies {
		urn := plugins.MaxComputeURN(dependency.Project, dependency.Dataset, dependency.Name)
		upstreams = append(upstreams, &v1beta2.Resource{
			Urn:     urn,
			Name:    dependency.Name,
			Type:    "table",
			Service: maxcomputeService,
		})
	}
	return upstreams
}

func buildColumns(dataType datatype.DataType) []*v1beta2.Column {
	if dataType.ID() != datatype.STRUCT {
		return nil
	}

	structType, ok := dataType.(datatype.StructType)
	if !ok {
		return nil
	}

	var columns []*v1beta2.Column
	for _, field := range structType.Fields {
		column := &v1beta2.Column{
			Name:     field.Name,
			DataType: dataTypeToString(field.Type),
			Columns:  buildColumns(field.Type),
		}
		columns = append(columns, column)
	}
	return columns
}

func (e *Extractor) buildTableAttributesData(schemaName, tableType string, tableInfo *tableschema.TableSchema) map[string]interface{} {
	attributesData := map[string]interface{}{}

	attributesData[attributesDataProjectName] = e.config.ProjectName
	attributesData[attributesDataSchema] = schemaName
	attributesData[attributesDataType] = tableType

	rb := common.ResourceBuilder{ProjectName: e.config.ProjectName}
	attributesData[attributesDataResourceURL] = rb.Table(schemaName, tableInfo.TableName)

	if tableInfo.ViewText != "" {
		attributesData[attributesDataSQL] = tableInfo.ViewText
	}

	if tableInfo.Lifecycle != 0 {
		attributesData[attributesDataLifecycle] = tableInfo.Lifecycle
	}

	var partitionNames []interface{}
	if len(tableInfo.PartitionColumns) > 0 {
		partitionNames = make([]interface{}, 0, len(tableInfo.PartitionColumns))
		for _, column := range tableInfo.PartitionColumns {
			name := column.Name
			if column.GenerateExpression != nil {
				if refCol := extractPartitionReferenceColumn(column.GenerateExpression); refCol != "" {
					name = refCol
				}
			}
			partitionNames = append(partitionNames, name)
		}
		attributesData[attributesDataPartitionFields] = partitionNames
	}

	return attributesData
}

func getDDLStatement(tableSchema *tableschema.TableSchema, projectName, schemaName string) (string, error) {
	var ddl string
	var err error

	switch {
	case tableSchema.IsVirtualView || tableSchema.IsMaterializedView:
		ddl, err = tableSchema.ToViewSQLString(projectName, schemaName, true, true, false)
	case !tableSchema.IsExternal:
		ddl, err = tableSchema.ToSQLString(projectName, schemaName, true)
	default:
		ddl, err = tableSchema.ToExternalSQLString(projectName, schemaName, true, nil, nil)
	}

	if err != nil {
		return "", err
	}
	return ddl, nil
}

func buildColumnAttributesData(column *tableschema.Column) map[string]interface{} {
	attributesData := map[string]interface{}{}

	if column == nil {
		return attributesData
	}

	if column.Label != "" {
		attributesData[attributesDataLabel] = column.Label
	}

	return attributesData
}

func (e *Extractor) buildPreview(ctx context.Context, t *odps.Table, tSchema *tableschema.TableSchema) ([]string, *structpb.ListValue, error) {
	maxPreviewRows := e.config.MaxPreviewRows
	if maxPreviewRows <= 0 {
		return nil, nil, nil
	}

	previewFields, previewRows, err := e.client.GetTablePreview(ctx, "", t, maxPreviewRows)
	if err != nil {
		e.logger.Error("failed to preview table", "table", t.Name(), "error", err)
		return nil, nil, err
	}

	if e.config.MixValues {
		tempRows := make([]interface{}, len(previewRows.GetValues()))
		for i, val := range previewRows.GetValues() {
			tempRows[i] = val.AsInterface()
		}

		tempRows, err = e.mixValuesIfNeeded(tempRows, time.Time(tSchema.LastModifiedTime).Unix())
		if err != nil {
			return nil, nil, fmt.Errorf("mix values: %w", err)
		}

		previewRows, err = structpb.NewList(tempRows)
		if err != nil {
			return nil, nil, fmt.Errorf("create preview list: %w", err)
		}
	}

	previewFields, previewRows = flattenJSONColumns(previewFields, previewRows, stringColumnSet(tSchema))

	return previewFields, previewRows, nil
}

func stringColumnSet(tableSchema *tableschema.TableSchema) map[string]bool {
	if tableSchema == nil {
		return nil
	}

	stringCols := make(map[string]bool, len(tableSchema.Columns))
	for _, col := range tableSchema.Columns {
		if dataTypeToString(col.Type) == datatype.STRING.String() {
			stringCols[col.Name] = true
		}
	}

	return stringCols
}

func flattenJSONColumns(fields []string, rows *structpb.ListValue, isStringCol map[string]bool) ([]string, *structpb.ListValue) {
	if rows == nil || len(isStringCol) == 0 {
		return fields, rows
	}

	subKeysByColumn := jsonSubKeys(fields, rows, isStringCol)
	if len(subKeysByColumn) == 0 {
		return fields, rows
	}

	newFields := make([]string, 0, len(fields))
	for colIdx, fieldName := range fields {
		newFields = append(newFields, fieldName)
		for _, subKey := range subKeysByColumn[colIdx] {
			newFields = append(newFields, fieldName+"."+subKey)
		}
	}

	newRows := make([]*structpb.Value, len(rows.GetValues()))
	for rowIdx, row := range rows.GetValues() {
		newRows[rowIdx] = expandJSONRow(row, len(newFields), subKeysByColumn)
	}

	return newFields, &structpb.ListValue{Values: newRows}
}

func jsonSubKeys(fields []string, rows *structpb.ListValue, isStringCol map[string]bool) map[int][]string {
	subKeysByColumn := map[int][]string{}
	for colIdx, fieldName := range fields {
		if !isStringCol[fieldName] {
			continue
		}

		keySet := map[string]struct{}{}
		for _, row := range rows.GetValues() {
			for key := range jsonObjectAt(row, colIdx) {
				keySet[key] = struct{}{}
			}
		}

		if len(keySet) == 0 {
			continue
		}

		keys := make([]string, 0, len(keySet))
		for key := range keySet {
			keys = append(keys, key)
		}

		sort.Strings(keys)
		subKeysByColumn[colIdx] = keys
	}

	return subKeysByColumn
}

func expandJSONRow(row *structpb.Value, width int, subKeysByColumn map[int][]string) *structpb.Value {
	expandedCells := make([]*structpb.Value, 0, width)
	for colIdx, cell := range row.GetListValue().GetValues() {
		expandedCells = append(expandedCells, cell)
		subKeys, ok := subKeysByColumn[colIdx]
		if !ok {
			continue
		}

		jsonObj := jsonObjectAt(row, colIdx)
		for _, subKey := range subKeys {
			expandedCells = append(expandedCells, structpb.NewStringValue(jsonText(jsonObj[subKey])))
		}
	}

	return structpb.NewListValue(&structpb.ListValue{Values: expandedCells})
}

func jsonObjectAt(row *structpb.Value, col int) map[string]json.RawMessage {
	cells := row.GetListValue().GetValues()
	if col >= len(cells) {
		return nil
	}

	text := strings.TrimSpace(cells[col].GetStringValue())
	if text == "" || text[0] != '{' {
		return nil
	}

	var obj map[string]json.RawMessage
	if json.Unmarshal([]byte(text), &obj) != nil {
		return nil
	}

	return obj
}

func jsonText(raw json.RawMessage) string {
	if len(raw) == 0 || string(raw) == "null" {
		return ""
	}

	var text string
	if json.Unmarshal(raw, &text) == nil {
		return text
	}

	return string(raw)
}

func (e *Extractor) mixValuesIfNeeded(rows []interface{}, rndSeed int64) ([]interface{}, error) {
	if !e.config.MixValues || len(rows) < 2 {
		return rows, nil
	}

	var table [][]any
	for _, row := range rows {
		arr, ok := row.([]any)
		if !ok {
			return nil, fmt.Errorf("row %d is not a slice", row)
		}
		table = append(table, arr)
	}

	numRows := len(table)
	numColumns := len(table[0])

	rndGen := e.randFn(rndSeed)
	for col := 0; col < numColumns; col++ {
		for row := 0; row < numRows; row++ {
			randomRow := rndGen(int64(numRows))

			table[row][col], table[randomRow][col] = table[randomRow][col], table[row][col]
		}
	}

	mixedRows := make([]any, numRows)
	for i, row := range table {
		mixedRows[i] = row
	}
	return mixedRows, nil
}

func dataTypeToString(dataType datatype.DataType) string {
	if dataType.ID() == datatype.MAP {
		return dataType.Name()
	}
	if dataType.ID() == datatype.ARRAY {
		return dataType.Name()
	}
	return dataType.ID().String()
}

func extractPartitionReferenceColumn(expr interface{ String() string }) string {
	if expr == nil {
		return ""
	}
	raw := expr.String()
	if raw == "" {
		return ""
	}

	// TruncTime.String() in the SDK returns: trunc_time(`columnName`, 'datePart')
	truncTimeRefRegex := regexp.MustCompile("(?i)trunc_time\\s*\\(\\s*`?([a-zA-Z_][a-zA-Z0-9_]*)`?")
	if matches := truncTimeRefRegex.FindStringSubmatch(raw); len(matches) == 2 {
		return matches[1]
	}

	return ""
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func init() {
	if err := registry.Extractors.Register(maxcomputeService, func() plugins.Extractor {
		return New(plugins.GetLog(), CreateClient, seededRandom)
	}); err != nil {
		panic(err)
	}
}

func seededRandom(seed int64) func(max int64) int64 {
	rnd := rand.New(rand.NewSource(seed)) //nolint:gosec
	return func(max int64) int64 {
		return rnd.Int63n(max)
	}
}

func CreateClient(_ context.Context, _ log.Logger, conf config.Config) (Client, error) {
	return client.New(conf)
}
