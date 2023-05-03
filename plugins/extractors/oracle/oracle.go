package oracle

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/goto/meteor/models"
	v1beta2 "github.com/goto/meteor/models/gotocompany/assets/v1beta2"
	"github.com/goto/meteor/plugins"
	"github.com/goto/meteor/plugins/sqlutil"
	"github.com/goto/meteor/registry"
	"github.com/goto/salt/log"
	"github.com/pkg/errors"
	_ "github.com/sijms/go-ora/v2"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
)

var summary string

// Config holds the set of configuration options for the extractor
type Config struct {
	ConnectionURL string `mapstructure:"connection_url" validate:"required"`
}

var sampleConfig = `connection_url: oracle://username:passwd@localhost:1521/xe`

var info = plugins.Info{
	Description:  "Table metadata oracle SQL Database.",
	SampleConfig: sampleConfig,
	Summary:      summary,
	Tags:         []string{"oss", "extractor"},
}

// Extractor manages the extraction of data from the extractor
type Extractor struct {
	plugins.BaseExtractor
	logger log.Logger
	config Config
	db     *sql.DB
}

// New returns a pointer to an initialized Extractor Object
func New(logger log.Logger) *Extractor {
	e := &Extractor{
		logger: logger,
	}
	e.BaseExtractor = plugins.NewBaseExtractor(info, &e.config)

	return e
}

// Init initializes the extractor
func (e *Extractor) Init(ctx context.Context, config plugins.Config) (err error) {
	if err = e.BaseExtractor.Init(ctx, config); err != nil {
		return err
	}

	// Create database connection
	e.db, err = connection(e.config)
	if err != nil {
		return errors.Wrap(err, "failed to create connection")
	}

	return
}

// Extract collects metadata from the source. Metadata is collected through the emitter
func (e *Extractor) Extract(ctx context.Context, emit plugins.Emit) (err error) {
	defer e.db.Close()

	// Get username
	userName, err := e.getUserName(e.db)
	if err != nil {
		e.logger.Error("failed to get the user name", "error", err)
		return
	}

	// Get DB name
	dbs, err := sqlutil.FetchDBs(e.db, e.logger, "select ora_database_name from dual")
	if err != nil {
		return
	}
	database := dbs[0]

	tableQuery := fmt.Sprintf(`SELECT object_name 
 		FROM all_objects
		WHERE object_type = 'TABLE'
		AND upper(owner) = upper('%s')`, userName)
	tables, err := sqlutil.FetchTablesInDB(e.db, database, tableQuery)
	if err != nil {
		e.logger.Error("failed to get tables, skipping database", "error", err)
		// continue
	}

	for _, table := range tables {
		result, err := e.getTableMetadata(e.db, database, table)
		if err != nil {
			e.logger.Error("failed to get table metadata, skipping table", "error", err)
			// continue
		}
		// Publish metadata to channel
		emit(models.NewRecord(result))
	}

	return nil
}

func (e *Extractor) getUserName(db *sql.DB) (userName string, err error) {
	sqlStr := `select user from dual`

	rows, err := db.Query(sqlStr)
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&userName)
		if err != nil {
			return
		}
	}
	return userName, err
}

// Prepares the list of tables and the attached metadata
func (e *Extractor) getTableMetadata(db *sql.DB, dbName string, tableName string) (result *v1beta2.Asset, err error) {
	var columns []*v1beta2.Column
	columns, err = e.getColumnMetadata(db, dbName, tableName)
	if err != nil {
		return result, nil
	}

	// get table row count
	sqlStr := `select count(*) from %s`
	rows, err := db.Query(fmt.Sprintf(sqlStr, tableName))
	if err != nil {
		err = fmt.Errorf("error running query: %w", err)
		return nil, err
	}
	var rowCount int64
	for rows.Next() {
		if err = rows.Scan(&rowCount); err != nil {
			e.logger.Error("failed to get fields", "error", err)
			continue
		}
	}
	table, err := anypb.New(&v1beta2.Table{
		Columns: columns,
		Profile: &v1beta2.TableProfile{
			TotalRows: rowCount,
		},
		Attributes: &structpb.Struct{}, // ensure attributes don't get overwritten if present
	})
	if err != nil {
		err = fmt.Errorf("error creating Any struct: %w", err)
		return nil, err
	}
	result = &v1beta2.Asset{
		Urn:     models.NewURN("oracle", e.UrnScope, "table", fmt.Sprintf("%s.%s", dbName, tableName)),
		Name:    tableName,
		Service: "Oracle",
		Type:    "table",
		Data:    table,
	}

	return
}

// Prepares the list of columns and the attached metadata
func (e *Extractor) getColumnMetadata(db *sql.DB, dbName string, tableName string) (result []*v1beta2.Column, err error) {
	sqlStr := `select utc.column_name, utc.data_type, 
			decode(utc.char_used, 'C', utc.char_length, utc.data_length) as data_length,
			utc.nullable, nvl(ucc.comments, '') as col_comment
			from USER_TAB_COLUMNS utc
			INNER JOIN USER_COL_COMMENTS ucc ON
			utc.column_name = ucc.column_name AND
			utc.table_name = ucc.table_name
			WHERE utc.table_name ='%s'`

	rows, err := db.Query(fmt.Sprintf(sqlStr, tableName))
	if err != nil {
		err = errors.Wrap(err, "failed to fetch data from query")
		return
	}

	for rows.Next() {
		var fieldName, dataType, isNullableString string
		var fieldDesc sql.NullString
		var length int
		if err = rows.Scan(&fieldName, &dataType, &length, &isNullableString, &fieldDesc); err != nil {
			e.logger.Error("failed to get fields", "error", err)
			continue
		}

		result = append(result, &v1beta2.Column{
			Name:        fieldName,
			DataType:    dataType,
			Description: fieldDesc.String,
			IsNullable:  isNullable(isNullableString),
			Length:      int64(length),
		})
	}
	return result, nil
}

// Convert nullable string to a boolean
func isNullable(value string) bool {
	return value == "Y"
}

// connection generates a connection string
func connection(cfg Config) (db *sql.DB, err error) {
	return sql.Open("oracle", cfg.ConnectionURL)
}

// Register the extractor to catalog
func init() {
	if err := registry.Extractors.Register("oracle", func() plugins.Extractor {
		return New(plugins.GetLog())
	}); err != nil {
		panic(err)
	}
}
