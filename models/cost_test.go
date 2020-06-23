package models_test

import (
	"testing"
	"time"

	"github.com/equinor/radix-cost-allocation/models"
	"github.com/stretchr/testify/assert"
)

func Test_cost_all_app_equal_requested(t *testing.T) {
	runs := getTestRuns()
	cost := models.NewCost(time.Now().Add(-1), time.Now(), runs)

	assert.Equal(t, cost.Applications[0].CostPercentage, cost.Applications[1].CostPercentage)
	assert.Equal(t, cost.Applications[1].CostPercentage, cost.Applications[2].CostPercentage)
	assert.Equal(t, cost.Applications[2].CostPercentage, cost.Applications[3].CostPercentage)

	assert.Equal(t, 0.25, cost.Applications[0].CostPercentage)
}

func Test_cost_one_app_double_requested(t *testing.T) {
	runs := getTestRuns()
	runs[0].Resources[0].Replicas = 4
	runs[1].Resources[0].Replicas = 4
	runs[2].Resources[0].Replicas = 4
	cost := models.NewCost(time.Now().Add(-1), time.Now(), runs)

	assert.Equal(t, 0.4, cost.GetCostBy("app-1").CostPercentage)
	assert.Equal(t, 0.2, cost.GetCostBy("app-2").CostPercentage)
	assert.Equal(t, 0.2, cost.GetCostBy("app-3").CostPercentage)
	assert.Equal(t, 0.2, cost.GetCostBy("app-4").CostPercentage)
}

func Test_cost_one_app_no_requested(t *testing.T) {
	runs := getTestRuns()
	runs[0].Resources[0].Replicas = 0
	runs[1].Resources[0].Replicas = 0
	runs[2].Resources[0].Replicas = 0
	cost := models.NewCost(time.Now().Add(-1), time.Now(), runs)

	oneThird := float64(1.0) / float64(3.0)
	assert.Equal(t, 0.0, cost.GetCostBy("app-1").CostPercentage)
	assert.Equal(t, oneThird, cost.GetCostBy("app-2").CostPercentage)
	assert.Equal(t, oneThird, cost.GetCostBy("app-3").CostPercentage)
	assert.Equal(t, oneThird, cost.GetCostBy("app-4").CostPercentage)
}

func getTestRuns() []models.Run {
	return []models.Run{
		{
			ID:                  1,
			ClusterCPUMillicore: 1000,
			Resources: []models.RequiredResources{
				{
					Application:     "app-1",
					Environment:     "env-1",
					Component:       "comp-1",
					CPUMillicore:    100,
					MemoryMegaBytes: 100,
					Replicas:        2,
				},
				{
					Application:     "app-2",
					Environment:     "env-1",
					Component:       "comp-1",
					CPUMillicore:    100,
					MemoryMegaBytes: 100,
					Replicas:        2,
				},
				{
					Application:     "app-3",
					Environment:     "env-1",
					Component:       "comp-1",
					CPUMillicore:    100,
					MemoryMegaBytes: 100,
					Replicas:        2,
				},
				{
					Application:     "app-4",
					Environment:     "env-1",
					Component:       "comp-1",
					CPUMillicore:    100,
					MemoryMegaBytes: 100,
					Replicas:        2,
				},
			},
		},
		{
			ID:                  2,
			ClusterCPUMillicore: 1000,
			Resources: []models.RequiredResources{
				{
					Application:     "app-1",
					Environment:     "env-1",
					Component:       "comp-1",
					CPUMillicore:    100,
					MemoryMegaBytes: 100,
					Replicas:        2,
				},
				{
					Application:     "app-2",
					Environment:     "env-1",
					Component:       "comp-1",
					CPUMillicore:    100,
					MemoryMegaBytes: 100,
					Replicas:        2,
				},
				{
					Application:     "app-3",
					Environment:     "env-1",
					Component:       "comp-1",
					CPUMillicore:    100,
					MemoryMegaBytes: 100,
					Replicas:        2,
				},
				{
					Application:     "app-4",
					Environment:     "env-1",
					Component:       "comp-1",
					CPUMillicore:    100,
					MemoryMegaBytes: 100,
					Replicas:        2,
				},
			},
		},
		{
			ID:                  3,
			ClusterCPUMillicore: 2000,
			Resources: []models.RequiredResources{
				{
					Application:     "app-1",
					Environment:     "env-1",
					Component:       "comp-1",
					CPUMillicore:    100,
					MemoryMegaBytes: 100,
					Replicas:        2,
				},
				{
					Application:     "app-2",
					Environment:     "env-1",
					Component:       "comp-1",
					CPUMillicore:    100,
					MemoryMegaBytes: 100,
					Replicas:        2,
				},
				{
					Application:     "app-3",
					Environment:     "env-1",
					Component:       "comp-1",
					CPUMillicore:    100,
					MemoryMegaBytes: 100,
					Replicas:        2,
				},
				{
					Application:     "app-4",
					Environment:     "env-1",
					Component:       "comp-1",
					CPUMillicore:    100,
					MemoryMegaBytes: 100,
					Replicas:        2,
				},
			},
		},
	}
}
