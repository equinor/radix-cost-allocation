package models

import (
	"time"
)

// Cost holds cost information between from and to time
type Cost struct {
	From         time.Time
	To           time.Time
	Applications []ApplicationCost
	runs         []Run
}

// ApplicationCost holds cost information with resources
type ApplicationCost struct {
	Name                   string
	WBS                    string
	CostPercentageByCPU    float64
	CostPercentageByMemory float64
}

// NewCost aggregate cost over a time period for applications
func NewCost(from, to time.Time, runs []Run) Cost {
	cost := Cost{
		From:         from,
		To:           to,
		Applications: aggregateCostBetweenDatesOnApplications(runs),
		runs:         runs,
	}
	return cost
}

// GetCostBy returns application by appName
func (cost Cost) GetCostBy(appName string) *ApplicationCost {
	for _, app := range cost.Applications {
		if app.Name == appName {
			return &app
		}
	}
	return nil
}

// aggregateCostBetweenDatesOnApplications calculates cost for an application
func aggregateCostBetweenDatesOnApplications(runs []Run) []ApplicationCost {
	totalRequestedCPU := totalRequestedCPU(runs)
	totalRequestedMemory := totalRequestedMemoryMegaBytes(runs)
	cpuPercentages := map[string]float64{}
	memoryPercentage := map[string]float64{}

	for _, runs := range runs {
		applications := runs.GetApplicationsRequiredResource()
		for _, application := range applications {
			cpuPercentages[application.Name] += runs.CPUWeightInPeriod(totalRequestedCPU) * application.RequestedCPUPercentageOfRun
			memoryPercentage[application.Name] += runs.MemoryWeightInPeriod(totalRequestedMemory) * application.RequestedMemoryPercentageOfRun
		}
	}

	applications := []ApplicationCost{}
	for appName, cpu := range cpuPercentages {
		applications = append(applications, ApplicationCost{
			Name:                   appName,
			CostPercentageByCPU:    cpu,
			CostPercentageByMemory: memoryPercentage[appName],
		})
	}
	return applications
}

func totalRequestedMemoryMegaBytes(runs []Run) int {
	memory := 0
	for _, run := range runs {
		memory += run.ClusterMemoryMegaByte
	}
	return memory
}

// TotalRequestedCPU total requested cpu for runs between from and to datetime
func totalRequestedCPU(runs []Run) int {
	cpu := 0
	for _, run := range runs {
		cpu += run.ClusterCPUMillicore
	}
	return cpu
}
