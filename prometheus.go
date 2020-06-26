package main

import (
	"context"
	"fmt"
	"time"

	"github.com/equinor/radix-cost-allocation/models"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

// PrometheusClient used to perform queries towards prometheus
type PrometheusClient struct {
	Address string
}

// GetRequiredResources get required resources for a given time
func (client PrometheusClient) GetRequiredResources(measuredTime time.Time) ([]models.RequiredResources, error) {
	requiredCPU, err := client.getRequiredCPUFromPrometheus(measuredTime)
	if err != nil {
		return nil, fmt.Errorf("Error getting cpu: %v", err)
	}
	requiredMemory, err := client.getRequiredMemoryFromPrometheus(measuredTime)
	if err != nil {
		return nil, fmt.Errorf("Error getting memory: %v", err)
	}
	requiredReplicas, err := client.getRequiredReplicasFromPrometheus(measuredTime)
	if err != nil {
		return nil, fmt.Errorf("Error getting replicas: %v", err)
	}
	reqResources, err := mapToRequiredResources(requiredCPU, requiredMemory, requiredReplicas)
	if err != nil {
		return nil, fmt.Errorf("Error mapping resources: %v", err)
	}
	return reqResources, nil
}

// GetClusterTotalCPUCoresFromPrometheus gets the number of cpu for a given time
func (client PrometheusClient) GetClusterTotalCPUCoresFromPrometheus(measuredTime time.Time) (int, error) {
	vector, err := client.getVectorPrometheus(measuredTime, "sum(instance:node_num_cpu:sum)")
	if err != nil {
		return 0, err
	}

	for _, val := range vector {
		return int(val.Value), nil
	}
	return 0, nil
}

// GetClusterTotalMemoryBytesFromPrometheus gets the sum of node_memory_MemTotal_bytes for a given time
func (client PrometheusClient) GetClusterTotalMemoryBytesFromPrometheus(measuredTime time.Time) (int, error) {
	vector, err := client.getVectorPrometheus(measuredTime, "sum(node_memory_MemTotal_bytes)")
	if err != nil {
		return 0, err
	}

	for _, val := range vector {
		return int(val.Value), nil
	}
	return 0, nil
}

func (client PrometheusClient) getRequiredCPUFromPrometheus(measuredTime time.Time) (model.Vector, error) {
	return client.getVectorPrometheus(measuredTime, "(sum(radix_operator_requested_cpu) by (application, environment, component))")
}

func (client PrometheusClient) getRequiredMemoryFromPrometheus(measuredTime time.Time) (model.Vector, error) {
	return client.getVectorPrometheus(measuredTime, "(sum(radix_operator_requested_memory) by (application, environment, component))")
}

func (client PrometheusClient) getRequiredReplicasFromPrometheus(measuredTime time.Time) (model.Vector, error) {
	return client.getVectorPrometheus(measuredTime, "(sum(radix_operator_requested_replicas) by (application, environment, component))")
}

func (client PrometheusClient) getVectorPrometheus(measuredTime time.Time, query string) (model.Vector, error) {
	promClient, err := api.NewClient(api.Config{
		Address: client.Address,
	})

	if err != nil {
		return nil, fmt.Errorf("error creating client: %v", err)
	}

	v1api := v1.NewAPI(promClient)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := v1api.Query(ctx, query, measuredTime)
	if err != nil {
		return nil, fmt.Errorf("error querying Prometheus: %v", err)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}
	if result.Type() != model.ValVector {
		return nil, fmt.Errorf("query prom required resources returned wrong datatype %s", result.Type())
	}
	vector := result.(model.Vector)

	return vector, nil
}

func mapToRequiredResources(requiredCPU, requiredMemory, requiredReplicas model.Vector) ([]models.RequiredResources, error) {
	reqResources := make(map[string]map[string]map[string]*models.RequiredResources)
	for _, replica := range requiredReplicas {
		metric := replica.Metric
		application, environment, component, wbs := getPropFromMetric(metric)

		req := models.RequiredResources{
			Application: application,
			Environment: environment,
			Component:   component,
			WBS:         wbs,
			Replicas:    int(replica.Value),
		}

		if reqResources[application] == nil {
			reqResources[application] = make(map[string]map[string]*models.RequiredResources)
		}
		if reqResources[application][environment] == nil {
			reqResources[application][environment] = make(map[string]*models.RequiredResources)
		}
		reqResources[application][environment][component] = &req
	}

	for _, cpu := range requiredCPU {
		metric := cpu.Metric
		application, environment, component, _ := getPropFromMetric(metric)

		req := reqResources[application][environment][component]
		req.CPUMillicore = int(cpu.Value)
	}

	for _, memory := range requiredMemory {
		metric := memory.Metric
		application, environment, component, _ := getPropFromMetric(metric)

		req := reqResources[application][environment][component]
		req.MemoryMegaBytes = int(memory.Value)
	}

	result := []models.RequiredResources{}
	for _, app := range reqResources {
		for _, env := range app {
			for _, comp := range env {
				result = append(result, *comp)
			}
		}
	}
	return result, nil
}

func getPropFromMetric(metric model.Metric) (app, env, component, wbs string) {
	if val, ok := metric["application"]; ok {
		app = string(val)
	}
	if val, ok := metric["environment"]; ok {
		env = string(val)
	}
	if val, ok := metric["component"]; ok {
		component = string(val)
	}
	if val, ok := metric["wbs"]; ok {
		wbs = string(val)
	}
	return
}
