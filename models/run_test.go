package models_test

import (
	"testing"

	"github.com/equinor/radix-cost-allocation/models"
	"github.com/stretchr/testify/assert"
)

func getTestRun() models.Run {
	return models.Run{
		ID:                    1,
		ClusterCPUMillicore:   1000,
		ClusterMemoryMegaByte: 1000,
	}
}

func Test_run_cpu_weight_in_period(t *testing.T) {
	run := getTestRun()

	run.ClusterCPUMillicore = 1000
	assert.Equal(t, 0.2, run.CPUWeightInPeriod(5000))
}

func Test_run_memory_weight_in_period(t *testing.T) {
	run := getTestRun()

	run.ClusterMemoryMegaByte = 1000
	assert.Equal(t, 0.2, run.MemoryWeightInPeriod(5000))
}

func Test_run_total_cpu_no_resources(t *testing.T) {
	run := getTestRun()
	run.Resources = []models.RequiredResources{}

	totalRequestedCPUByApp := run.RequestedCPUByApplications()
	assert.Equal(t, 0, totalRequestedCPUByApp)
}

func Test_run_total_cpu_requested(t *testing.T) {
	run := getTestRun()
	run.Resources = []models.RequiredResources{
		{
			CPUMillicore: 100,
			Replicas:     2,
		},
		{
			CPUMillicore: 100,
			Replicas:     1,
		},
		{
			CPUMillicore: 100,
			Replicas:     1,
		},
	}

	totalRequestedCPUByApp := run.RequestedCPUByApplications()
	assert.Equal(t, 400, totalRequestedCPUByApp)
}

func Test_run_total_memory_no_resources(t *testing.T) {
	run := getTestRun()
	run.Resources = []models.RequiredResources{}

	totalRequestedMemoryByApp := run.RequestedMemoryByApplications()
	assert.Equal(t, 0, totalRequestedMemoryByApp)
}

func Test_run_total_memory_requested(t *testing.T) {
	run := getTestRun()
	run.Resources = []models.RequiredResources{
		{
			MemoryMegaBytes: 100,
			Replicas:        2,
		},
		{
			MemoryMegaBytes: 100,
			Replicas:        1,
		},
		{
			MemoryMegaBytes: 100,
			Replicas:        1,
		},
	}

	totalRequestedMemoryByApp := run.RequestedMemoryByApplications()
	assert.Equal(t, 400, totalRequestedMemoryByApp)
}

func Test_get_application_requested_resources(t *testing.T) {
	run := loadDefaultResources(getTestRun())
	applications := run.GetApplicationsRequiredResource()

	assert.Equal(t, 2, len(applications))

	for _, app := range applications {
		if app.Name == "app-1" {
			assert.Equal(t, 600, app.RequestedCPUMillicore)
			assert.Equal(t, 300, app.RequestedMemoryMegaByte)
			assert.Equal(t, 0.6, app.RequestedCPUPercentageOfRun)
			assert.Equal(t, 0.3, app.RequestedMemoryPercentageOfRun)
		}
		if app.Name == "app-2" {
			assert.Equal(t, 400, app.RequestedCPUMillicore)
			assert.Equal(t, 700, app.RequestedMemoryMegaByte)
			assert.Equal(t, 0.4, app.RequestedCPUPercentageOfRun)
			assert.Equal(t, 0.7, app.RequestedMemoryPercentageOfRun)
		}
	}

}

func loadDefaultResources(run models.Run) models.Run {
	run.Resources = []models.RequiredResources{
		{
			Application:     "app-1",
			Environment:     "env-1",
			Component:       "comp-1",
			WBS:             "1",
			CPUMillicore:    200,
			MemoryMegaBytes: 100,
			Replicas:        2,
		},
		{
			Application:     "app-1",
			Environment:     "env-2",
			Component:       "comp-1",
			WBS:             "1",
			CPUMillicore:    200,
			MemoryMegaBytes: 100,
			Replicas:        1,
		},
		{
			Application:     "app-2",
			Environment:     "env-1",
			Component:       "comp-1",
			WBS:             "2",
			CPUMillicore:    100,
			MemoryMegaBytes: 100,
			Replicas:        1,
		},
		{
			Application:     "app-2",
			Environment:     "env-1",
			Component:       "comp-2",
			WBS:             "2",
			CPUMillicore:    100,
			MemoryMegaBytes: 200,
			Replicas:        1,
		},
		{
			Application:     "app-2",
			Environment:     "env-2",
			Component:       "comp-2",
			WBS:             "2",
			CPUMillicore:    100,
			MemoryMegaBytes: 200,
			Replicas:        2,
		},
	}
	return run
}
