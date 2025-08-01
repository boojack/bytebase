package catalog

// This file defines the database state during walk-through. The only two cases will use database state:
//   1. during walk-through
//   2. the underlying implementation of Finder

import (
	"fmt"
	"strings"

	storepb "github.com/bytebase/bytebase/backend/generated-go/store"
)

func newDatabaseState(d *storepb.DatabaseSchemaMetadata, context *FinderContext) *DatabaseState {
	database := &DatabaseState{
		ctx:          context.Copy(),
		name:         d.Name,
		characterSet: d.CharacterSet,
		collation:    d.Collation,
		dbType:       context.EngineType,
		schemaSet:    make(schemaStateMap),
		usable:       true,
	}

	for _, schema := range d.Schemas {
		database.schemaSet[schema.Name] = newSchemaState(schema, database.ctx)
	}

	for _, schema := range d.Schemas {
		for _, view := range schema.Views {
			for _, dependencyColumn := range view.DependencyColumns {
				if schemaState, exist := database.schemaSet[dependencyColumn.Schema]; exist {
					if tableState, exist := schemaState.getTable(dependencyColumn.Table); exist {
						tableState.dependencyView[fmt.Sprintf("%q.%q", schema.Name, view.Name)] = true
						if columnState, exist := tableState.columnSet[dependencyColumn.Column]; exist {
							columnState.dependencyView[fmt.Sprintf("%q.%q", schema.Name, view.Name)] = true
						}
					}
				}
			}
		}
	}

	return database
}

func newSchemaState(s *storepb.SchemaMetadata, context *FinderContext) *SchemaState {
	schema := &SchemaState{
		ctx:           context.Copy(),
		name:          s.Name,
		tableSet:      make(tableStateMap),
		viewSet:       make(viewStateMap),
		identifierMap: make(identifierMap),
	}

	for _, table := range s.Tables {
		tableState := newTableState(table, context)
		schema.tableSet[table.Name] = tableState

		schema.identifierMap[table.Name] = true
		for indexName := range tableState.indexSet {
			schema.identifierMap[indexName] = true
		}
	}

	for _, view := range s.Views {
		schema.viewSet[view.Name] = newViewState(view)

		schema.identifierMap[view.Name] = true
	}

	return schema
}

func newViewState(v *storepb.ViewMetadata) *ViewState {
	return &ViewState{
		name:       v.Name,
		definition: newStringPointer(v.Definition),
		comment:    newStringPointer(v.Comment),
	}
}

func newTableState(t *storepb.TableMetadata, context *FinderContext) *TableState {
	table := &TableState{
		name:           t.Name,
		engine:         newStringPointer(t.Engine),
		collation:      newStringPointer(t.Collation),
		comment:        newStringPointer(t.Comment),
		columnSet:      make(columnStateMap),
		indexSet:       make(IndexStateMap),
		dependencyView: make(map[string]bool),
	}

	for i, column := range t.Columns {
		columnName := column.Name
		switch context.EngineType {
		case storepb.Engine_MYSQL, storepb.Engine_TIDB, storepb.Engine_MARIADB, storepb.Engine_OCEANBASE:
			columnName = strings.ToLower(columnName)
		default:
			// For other engine types, keep the original column name without normalization
		}
		table.columnSet[columnName] = newColumnState(column, i+1)
	}

	for _, index := range t.Indexes {
		indexName := index.Name
		switch context.EngineType {
		case storepb.Engine_MYSQL, storepb.Engine_TIDB, storepb.Engine_MARIADB, storepb.Engine_OCEANBASE:
			indexName = strings.ToLower(indexName)
		default:
			// For other engine types, keep the original index name without normalization
		}
		table.indexSet[indexName] = newIndexState(index)
	}

	return table
}

func newColumnState(c *storepb.ColumnMetadata, position int) *ColumnState {
	defaultValue := (*string)(nil)
	if c.Default != "" {
		defaultValue = copyStringPointer(&c.Default)
	}
	return &ColumnState{
		name:           c.Name,
		position:       newIntPointer(position),
		defaultValue:   defaultValue,
		nullable:       newBoolPointer(c.Nullable),
		columnType:     newStringPointer(c.Type),
		characterSet:   newStringPointer(c.CharacterSet),
		collation:      newStringPointer(c.Collation),
		comment:        newStringPointer(c.Comment),
		dependencyView: make(map[string]bool),
	}
}

func newIndexState(i *storepb.IndexMetadata) *IndexState {
	index := &IndexState{
		name:           i.Name,
		indexType:      newStringPointer(i.Type),
		unique:         newBoolPointer(i.Unique),
		primary:        newBoolPointer(i.Primary),
		visible:        newBoolPointer(i.Visible),
		comment:        newStringPointer(i.Comment),
		expressionList: copyStringSlice(i.Expressions),
		// We rudely think that pk and uk are constraints here.
		// But in fact, we can create uk by CREATE UNIQUE INDEX statements.
		isConstraint: *newBoolPointer(i.Primary || i.Unique),
	}
	return index
}

// DatabaseState is the state for walk-through.
type DatabaseState struct {
	ctx          *FinderContext
	name         string
	characterSet string
	collation    string
	dbType       storepb.Engine
	schemaSet    schemaStateMap
	deleted      bool
	usable       bool
}

type TableIndexFind struct {
	SchemaName string
	TableName  string
}

// Index returns the index map of the table.
func (d *DatabaseState) Index(tableIndexFind *TableIndexFind) *IndexStateMap {
	schema, exists := d.schemaSet[tableIndexFind.SchemaName]
	if !exists {
		return nil
	}
	return schema.Index(tableIndexFind)
}

// Usable returns the usable of the database state.
func (d *DatabaseState) Usable() bool {
	return d.usable
}

// HasNoTable returns true if the current database has no table.
func (d *DatabaseState) HasNoTable() bool {
	for _, schema := range d.schemaSet {
		if len(schema.tableSet) != 0 {
			return false
		}
	}

	return true
}

func (d *DatabaseState) HasSchema(name string) bool {
	_, exists := d.schemaSet[name]
	return exists
}

// DatabaseName returns the database name.
func (d *DatabaseState) DatabaseName() string {
	return d.name
}

// IndexFind is for find index.
type IndexFind struct {
	SchemaName string
	TableName  string
	IndexName  string
}

// FindIndex finds the index.
func (d *DatabaseState) FindIndex(find *IndexFind) (string, *IndexState) {
	switch d.dbType {
	case storepb.Engine_MYSQL, storepb.Engine_TIDB, storepb.Engine_MARIADB, storepb.Engine_OCEANBASE:
		find.IndexName = strings.ToLower(find.IndexName)
	default:
		// For other engine types, keep the original index name without normalization
	}
	// There are two cases to find a index:
	// 1. find an index in specific table. e.g. MySQL and TiDB.
	// 2. find an index in the schema. e.g. PostgreSQL.
	// In PostgreSQL, the index name is unique in a schema, not a table.
	// In MySQL and TiDB, the index name is unique in a table.
	// So for case one, we need match table name, but for case two, we don't need.
	needMatchTable := (d.dbType != storepb.Engine_POSTGRES || find.SchemaName == "" || find.TableName != "")
	if needMatchTable {
		schema, exists := d.schemaSet[find.SchemaName]
		if !exists {
			return "", nil
		}
		table, exists := schema.getTable(find.TableName)
		if !exists {
			return "", nil
		}
		index, exists := table.indexSet[find.IndexName]
		if !exists {
			return "", nil
		}
		return table.name, index
	}
	for _, schema := range d.schemaSet {
		if schema.name != find.SchemaName {
			continue
		}
		for _, table := range schema.tableSet {
			// no need to further match table name because index is already unique in the schema
			index, exists := table.indexSet[find.IndexName]
			if !exists {
				return "", nil
			}
			return table.name, index
		}
	}
	return "", nil
}

// PrimaryKeyFind is for finding primary key.
type PrimaryKeyFind struct {
	SchemaName string
	TableName  string
}

// FindPrimaryKey finds the primary key.
func (d *DatabaseState) FindPrimaryKey(find *PrimaryKeyFind) *IndexState {
	for _, schema := range d.schemaSet {
		if schema.name != find.SchemaName {
			continue
		}
		for _, table := range schema.tableSet {
			if !compareIdentifier(table.name, find.TableName, schema.ctx.IgnoreCaseSensitive) {
				continue
			}
			for _, index := range table.indexSet {
				if index.primary != nil && *index.primary {
					return index
				}
			}
		}
	}
	return nil
}

// ColumnCount is for counting columns.
type ColumnCount struct {
	SchemaName string
	TableName  string
	ColumnType string
}

// CountColumn counts columns.
func (d *DatabaseState) CountColumn(count *ColumnCount) int {
	schema, exists := d.schemaSet[count.SchemaName]
	if !exists {
		return 0
	}
	table, exists := schema.getTable(count.TableName)
	if !exists {
		return 0
	}
	res := 0
	for _, column := range table.columnSet {
		if column.columnType != nil && strings.EqualFold(*column.columnType, count.ColumnType) {
			res++
		}
	}
	return res
}

// ColumnFind is for finding column.
type ColumnFind struct {
	SchemaName string
	TableName  string
	ColumnName string
}

// FindColumn finds the column.
func (d *DatabaseState) FindColumn(find *ColumnFind) *ColumnState {
	switch d.dbType {
	case storepb.Engine_MYSQL, storepb.Engine_TIDB, storepb.Engine_MARIADB, storepb.Engine_OCEANBASE:
		find.ColumnName = strings.ToLower(find.ColumnName)
	default:
		// For other engine types, keep the original column name without normalization
	}
	schema, exists := d.schemaSet[find.SchemaName]
	if !exists {
		return nil
	}
	table, exists := schema.getTable(find.TableName)
	if !exists {
		return nil
	}
	column, exists := table.columnSet[find.ColumnName]
	if !exists {
		return nil
	}
	return column
}

// TableFind is for find table.
type TableFind struct {
	SchemaName string
	TableName  string
}

// FindTable finds the table.
func (d *DatabaseState) FindTable(find *TableFind) *TableState {
	schema, exists := d.schemaSet[find.SchemaName]
	if !exists {
		return nil
	}
	table, exists := schema.getTable(find.TableName)
	if !exists {
		return nil
	}
	return table
}

type identifierMap map[string]bool

// SchemaState is the state for walk-through.
type SchemaState struct {
	ctx      *FinderContext
	name     string
	tableSet tableStateMap
	viewSet  viewStateMap

	// PostgreSQL specific fields
	// All relation names in PostgreSQL must be distinct in schema level.
	identifierMap identifierMap
}

func (s *SchemaState) Index(tableIndexFind *TableIndexFind) *IndexStateMap {
	table, exists := s.getTable(tableIndexFind.TableName)
	if !exists {
		return nil
	}
	return table.Index(tableIndexFind)
}

type schemaStateMap map[string]*SchemaState

// TableState is the state for walk-through.
type TableState struct {
	name string
	// engine isn't supported for Postgres, Snowflake, SQLite.
	engine *string
	// collation isn't supported for Postgres, ClickHouse, Snowflake, SQLite.
	collation *string
	// comment isn't supported for SQLite.
	comment   *string
	columnSet columnStateMap
	// indexSet isn't supported for ClickHouse, Snowflake.
	indexSet IndexStateMap

	// dependencyView is used to record the dependency view for the table.
	// Used to check if the table is used by any view.
	dependencyView map[string]bool
}

// ListColumns return the list of columns.
func (t *TableState) ListColumns() []*ColumnState {
	var columns []*ColumnState
	for _, v := range t.columnSet {
		columns = append(columns, v)
	}
	return columns
}

// CountIndex return the index total number.
func (t *TableState) CountIndex() int {
	return len(t.indexSet)
}

// Index return the index map of table.
func (t *TableState) Index(_ *TableIndexFind) *IndexStateMap {
	return &t.indexSet
}

func (t *TableState) copy() *TableState {
	return &TableState{
		name:      t.name,
		engine:    copyStringPointer(t.engine),
		collation: copyStringPointer(t.collation),
		comment:   copyStringPointer(t.comment),
		columnSet: t.columnSet.copy(),
		indexSet:  t.indexSet.copy(),
	}
}

type tableStateMap map[string]*TableState

// IndexState is the state for walk-through.
type IndexState struct {
	name string
	// This could refer to a column or an expression.
	expressionList []string
	// Type isn't supported for SQLite.
	indexType *string
	unique    *bool
	primary   *bool
	// Visible isn't supported for Postgres, SQLite.
	visible *bool
	// Comment isn't supported for SQLite.
	comment *string

	// PostgreSQL specific fields.

	// PostgreSQL treats INDEX and CONSTRAINT differently.
	isConstraint bool
}

func (idx *IndexState) copy() *IndexState {
	return &IndexState{
		name:           idx.name,
		expressionList: copyStringSlice(idx.expressionList),
		indexType:      copyStringPointer(idx.indexType),
		unique:         copyBoolPointer(idx.unique),
		primary:        copyBoolPointer(idx.primary),
		visible:        copyBoolPointer(idx.visible),
		comment:        copyStringPointer(idx.comment),
		isConstraint:   idx.isConstraint,
	}
}

// Unique returns the unique for the index.
func (idx *IndexState) Unique() bool {
	if idx.unique != nil {
		return *idx.unique
	}
	return false
}

// Primary returns the primary for the index.
func (idx *IndexState) Primary() bool {
	if idx.primary != nil {
		return *idx.primary
	}
	return false
}

// ExpressionList returns the expression list for the index.
func (idx *IndexState) ExpressionList() []string {
	return idx.expressionList
}

type IndexStateMap map[string]*IndexState

func (m IndexStateMap) copy() IndexStateMap {
	res := make(IndexStateMap)
	for k, v := range m {
		res[k] = v.copy()
	}
	return res
}

// ColumnState is the state for walk-through.
type ColumnState struct {
	name         string
	position     *int
	defaultValue *string
	// nullable isn't supported for ClickHouse.
	nullable   *bool
	columnType *string
	// characterSet isn't supported for Postgres, ClickHouse, SQLite.
	characterSet *string
	// collation isn't supported for ClickHouse, SQLite.
	collation *string
	// comment isn't supported for SQLite.
	comment *string

	// dependencyView is used to record the dependency view for the column.
	// Used to check if the column is used by any view.
	dependencyView map[string]bool
}

func (col *ColumnState) copy() *ColumnState {
	return &ColumnState{
		name:         col.name,
		position:     copyIntPointer(col.position),
		defaultValue: copyStringPointer(col.defaultValue),
		nullable:     copyBoolPointer(col.nullable),
		columnType:   copyStringPointer(col.columnType),
		characterSet: copyStringPointer(col.characterSet),
		collation:    copyStringPointer(col.collation),
		comment:      copyStringPointer(col.comment),
	}
}

// Nullable returns nullable for the column.
func (col *ColumnState) Nullable() bool {
	return col.nullable != nil && *col.nullable
}

// Type returns type for the column.
func (col *ColumnState) Type() string {
	if col.columnType != nil {
		return *col.columnType
	}
	return ""
}

// HasDefault returns if column has default value.
func (col *ColumnState) HasDefault() bool {
	switch strings.ToLower(col.Type()) {
	case "serial", "smallserial", "bigserial":
		return true
	default:
		// For other column types, check if defaultValue is not nil
	}
	return col.defaultValue != nil
}

type columnStateMap map[string]*ColumnState

func (m columnStateMap) copy() columnStateMap {
	res := make(columnStateMap)
	for k, v := range m {
		res[k] = v.copy()
	}
	return res
}

// ViewState is the state for walk-through.
type ViewState struct {
	name       string
	definition *string
	comment    *string
}
type viewStateMap map[string]*ViewState

func copyStringPointer(p *string) *string {
	if p != nil {
		v := *p
		return &v
	}
	return nil
}

func copyBoolPointer(p *bool) *bool {
	if p != nil {
		v := *p
		return &v
	}
	return nil
}

func copyIntPointer(p *int) *int {
	if p != nil {
		v := *p
		return &v
	}
	return nil
}

func copyStringSlice(in []string) []string {
	var res []string
	res = append(res, in...)
	return res
}

func newEmptyStringPointer() *string {
	res := ""
	return &res
}

func newStringPointer(v string) *string {
	return &v
}

func newIntPointer(v int) *int {
	return &v
}

func newTruePointer() *bool {
	v := true
	return &v
}

func newFalsePointer() *bool {
	v := false
	return &v
}

func newBoolPointer(v bool) *bool {
	return &v
}
