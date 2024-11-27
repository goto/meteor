package tableau

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/goto/meteor/models"
	v1beta2 "github.com/goto/meteor/models/gotocompany/assets/v1beta2"
	"github.com/goto/meteor/plugins"
)

const (
	bigquery                = "bigquery"
	mssql                   = "mssql"
	maxcompute              = "maxcompute"
	maxcomputeDefaultSchema = "default"
)

// https://help.tableau.com/current/api/rest_api/en-us/REST/rest_api_concepts_connectiontype.htm
var connectionTypeMap = map[string]string{
	"sqlserver":       mssql,
	"maxcompute_jdbc": maxcompute,
}

func mapConnectionTypeToSource(ct string) string {
	s, ok := connectionTypeMap[ct]
	if !ok {
		return ct
	}
	return s
}

type Pagination struct {
	PageNumber     string `json:"pageNumber"`
	PageSize       string `json:"pageSize"`
	TotalAvailable string `json:"totalAvailable"`
}

type Project struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Table struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Schema   string   `json:"schema"`
	FullName string   `json:"fullName"`
	Database Database `json:"database"`
}

type Owner struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Workbook struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	ProjectName    string    `json:"projectName"`
	URI            string    `json:"uri"`
	Description    string    `json:"description"`
	Owner          Owner     `json:"owner"`
	Sheets         []*Sheet  `json:"sheets"`
	UpstreamTables []*Table  `json:"upstreamTables"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type Sheet struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// https://help.tableau.com/current/api/metadata_api/en-us/docs/meta_api_model.html
type DatabaseInterface interface {
	CreateResource(tableInfo Table) (resource *v1beta2.Resource)
}

type Database map[string]interface{}

type DatabaseServer struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	ConnectionType string `json:"connectionType"`
	Description    string `json:"description"`
	HostName       string `json:"hostName"`
	Port           int    `json:"port"`
	Service        string `json:"service"`
}

// parseBigQueryTableFullName would parse table full name into splitted strings (project id, dataset, table name)
// Full name found in tableau with BiqQuery source is like this
// sometimes table name can also be the same as full name (e.g. project-id.schema.table1)
// `project-id.schema.table1`
// [project-id.schema].[table1]`
func parseBigQueryTableFullName(fullName string) ([]string, bool) {
	omitedChars := "`" + "\\[\\]"
	re := regexp.MustCompile("[" + omitedChars + "]")
	cleanedFN := re.ReplaceAllString(fullName, "")
	splittedFN := strings.Split(cleanedFN, ".")
	return splittedFN, len(splittedFN) == 3
}

func (dbs *DatabaseServer) CreateResource(tableInfo Table) *v1beta2.Resource {
	source := mapConnectionTypeToSource(dbs.ConnectionType)

	var urn string
	switch source {
	case bigquery:
		// bigquery::sample-project/dataset_a/invoice
		// sometimes table name can be the same as full name (e.g. project-id.schema.table1), so we build URN with the full name instead
		fullNameSplitted, ok := parseBigQueryTableFullName(tableInfo.FullName)
		if !ok {
			// assume fullNameSplitted[0] is the project ID
			urn = plugins.BigQueryURN(fullNameSplitted[0], tableInfo.Schema, tableInfo.Name)
			break
		}
		urn = plugins.BigQueryURN(fullNameSplitted[0], fullNameSplitted[1], fullNameSplitted[2])
	case maxcompute:
		// depends on the JDBC MaxCompute driver version, tableInfo can consist another schema information or only can get the default schema
		if tableInfo.Schema != dbs.Name {
			urn = plugins.MaxComputeURN(dbs.Name, tableInfo.Schema, tableInfo.Name)
		} else {
			urn = plugins.MaxComputeURN(dbs.Name, maxcomputeDefaultSchema, tableInfo.Name)
		}
	default:
		// postgres::postgres:5432/postgres/user
		host := fmt.Sprintf("%s:%d", dbs.HostName, dbs.Port)
		urn = models.NewURN(source, host, "table", fmt.Sprintf("%s.%s", dbs.Name, tableInfo.Name))
	}
	return &v1beta2.Resource{
		Urn:     urn,
		Type:    "table",
		Service: source,
	}
}

type CloudFile struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	ConnectionType string `json:"connectionType"`
	Description    string `json:"description"`
	Provider       string `json:"provider"`
	FileExtension  string `json:"fileExtension"`
	MimeType       string `json:"mimeType"`
	RequestURL     string `json:"requestUrl"`
}

func (cf *CloudFile) CreateResource(tableInfo Table) *v1beta2.Resource {
	source := mapConnectionTypeToSource(cf.ConnectionType)
	urn := models.NewURN(source, cf.Provider, "bucket", fmt.Sprintf("%s/%s", cf.Name, tableInfo.Name))
	return &v1beta2.Resource{
		Urn:     urn,
		Type:    "bucket", // TODO need to check what would be the appropriate type for this
		Service: source,
	}
}

type File struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	ConnectionType string `json:"connectionType"`
	Description    string `json:"description"`
	FilePath       string `json:"filePath"`
}

func (f *File) CreateResource(tableInfo Table) *v1beta2.Resource {
	source := mapConnectionTypeToSource(f.ConnectionType)
	urn := models.NewURN(source, f.FilePath, "bucket", fmt.Sprintf("%s.%s", f.Name, tableInfo.Name))
	return &v1beta2.Resource{
		Urn:     urn,
		Type:    "bucket", // TODO need to check what would be the appropriate type for this
		Service: source,
	}
}

type WebDataConnector struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	ConnectionType string `json:"connectionType"`
	Description    string `json:"description"`
	ConnectorURL   string `json:"connectorUrl"`
}

func (wdc *WebDataConnector) CreateResource(tableInfo Table) *v1beta2.Resource {
	source := mapConnectionTypeToSource(wdc.ConnectionType)
	urn := models.NewURN(source, wdc.ConnectorURL, "table", fmt.Sprintf("%s.%s", wdc.Name, tableInfo.Name))
	return &v1beta2.Resource{
		Urn:     urn,
		Type:    "table", // TODO need to check what would be the appropriate type for this
		Service: source,
	}
}
