package repository

// NodeBulkDto defines an object with fields for bulk insert/update used by the Repository interface
// Fields in this struct is one-to-one mapped to the SQL Server cost.node_upsert_type table value type.
// Order of fields in the struct must match the order of fields in the SQL Server type
type NodeBulkDto struct {
	Name     string
	NodePool string
}
