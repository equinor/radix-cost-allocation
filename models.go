package main

import "time"

// Run holds all required resources for a time
type Run struct {
	ID              int64
	MeasuredTimeUTC time.Time
	Resources       []RequiredResources
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
