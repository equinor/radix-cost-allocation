package repository

// Repository defines methods for bulk upserting container and node resources to a database
type Repository interface {
	BulkUpsertContainers([]ContainerBulkTvp) error
	BulkUpsertNodes([]NodeBulkTvp) error
}
