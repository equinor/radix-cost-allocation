package repository

import "github.com/equinor/radix-cost-allocation/pkg/models"

type Repository interface {
	BulkUpsertContainers([]models.ContainerBulkTvp) error
	BulkUpsertNodes([]models.NodeBulkTvp) error
}
