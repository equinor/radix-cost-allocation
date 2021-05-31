package repository

// Repository defines methods for bulk upserting container and node resources to a database
type Repository interface {
	BulkUpsertContainers([]ContainerBulkDto) error
	BulkUpsertNodes([]NodeBulkDto) error
}
