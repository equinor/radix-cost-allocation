package models

import "time"

// Run holds all required resources for a time
type Run struct {
	ID                    int64
	MeasuredTimeUTC       time.Time
	ClusterCPUMillicore   int
	ClusterMemoryMegaByte int
	Resources             []RequiredResources
}

// RequiredResources holds required resources for a single component
type RequiredResources struct {
	ID              int64
	WBS             string
	Application     string
	Environment     string
	Component       string
	CPUMillicore    int
	MemoryMegaBytes int
	Replicas        int
}

type Application struct {
	Name                           string
	RequestedCPUMillicore          int
	RequestedMemoryMegaByte        int
	RequestedCPUPercentageOfRun    float64
	RequestedMemoryPercentageOfRun float64
}

// CPUWeightInPeriod weight of a run for a period
func (run Run) CPUWeightInPeriod(totalRequestedCPUForPeriod int) float64 {
	if totalRequestedCPUForPeriod == 0 {
		return 1
	}

	return float64(run.ClusterCPUMillicore) / float64(totalRequestedCPUForPeriod)
}

// MemoryWeightInPeriod weight of a run for a period
func (run Run) MemoryWeightInPeriod(totalRequestedMemoryForPeriod int) float64 {
	if totalRequestedMemoryForPeriod == 0 {
		return 1
	}

	return float64(run.ClusterMemoryMegaByte) / float64(totalRequestedMemoryForPeriod)
}

// RequestedCPUByApplications total requested cpu by applications for a run
func (run Run) RequestedCPUByApplications() int {
	cpu := 0
	for _, req := range run.Resources {
		cpu += req.CPUMillicore * req.Replicas
	}
	return cpu
}

// RequestedMemoryByApplications total requested memory by applications for a run
func (run Run) RequestedMemoryByApplications() int {
	memory := 0
	for _, req := range run.Resources {
		memory += req.MemoryMegaBytes * req.Replicas
	}
	return memory
}

// GetApplicationsRequiredResource returns resource requested aggregated on application
func (run Run) GetApplicationsRequiredResource() []Application {
	totalCPURequestedForRun := float64(run.RequestedCPUByApplications())
	totalMemoryRequestedForRun := float64(run.RequestedMemoryByApplications())

	requiredCPU := map[string]int{}
	requiredMemory := map[string]int{}
	for _, resource := range run.Resources {
		requiredCPU[resource.Application] += resource.CPUMillicore * resource.Replicas
		requiredMemory[resource.Application] += resource.MemoryMegaBytes * resource.Replicas
	}
	applications := []Application{}
	for key, val := range requiredCPU {
		applications = append(applications, Application{
			Name:                           key,
			RequestedCPUMillicore:          val,
			RequestedMemoryMegaByte:        requiredMemory[key],
			RequestedCPUPercentageOfRun:    float64(val) / totalCPURequestedForRun,
			RequestedMemoryPercentageOfRun: float64(requiredMemory[key]) / totalMemoryRequestedForRun,
		})
	}
	return applications
}
