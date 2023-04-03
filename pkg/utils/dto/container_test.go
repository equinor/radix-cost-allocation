package dto

import (
	"fmt"
	"testing"
	"time"

	"github.com/equinor/radix-cost-allocation/pkg/repository"
	"github.com/equinor/radix-cost-allocation/pkg/utils/clock"
	"github.com/equinor/radix-operator/pkg/apis/kube"
	radixv1 "github.com/equinor/radix-operator/pkg/apis/radix/v1"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestMapContainerBulkDtoFromPod(t *testing.T) {

	// Pod, correct label, container running, resources, rr exist
	t.Run("correct label, container running, resources set, rr exist", func(t *testing.T) {
		t.Parallel()
		podName, env, node, app, wbs, componentName, containerName, containerId, millicores, memory, startedAt, currentTime :=
			"pod1", "prod", "node1", "app1", "wbs", "comp", "c1", "cid1", int64(100), int64(1024), time.Date(2020, 1, 1, 1, 1, 1, 0, time.UTC), time.Date(2020, 2, 1, 1, 1, 1, 0, time.UTC)
		pod := buildPodForTest(podName, fmt.Sprintf("%s-%s", app, env), node,
			setPodAppLabel(app),
			setPodComponentLabel(componentName),
			appendPodContainer(
				buildContainerForTest(containerName,
					setContainerResourceRequest(corev1.ResourceCPU, *resource.NewMilliQuantity(millicores, resource.DecimalSI)),
					setContainerResourceRequest(corev1.ResourceMemory, *resource.NewQuantity(memory, resource.DecimalSI)),
				),
			),
			appendPodContainerStatus(
				buildContainerStatusForTest(containerName, containerId,
					setContainerStateRunning(startedAt),
				),
			),
		)
		expected := []repository.ContainerBulkDto{{
			ContainerID:          containerId,
			ContainerName:        containerName,
			PodName:              podName,
			ApplicationName:      app,
			EnvironmentName:      env,
			ComponentName:        componentName,
			Wbs:                  wbs,
			StartedAt:            startedAt,
			LastKnowRunningAt:    currentTime,
			CPURequestMillicores: millicores,
			MemoryRequestBytes:   memory,
			NodeName:             node,
		}}
		rrMap := map[string]*radixv1.RadixRegistration{"app1": {Spec: radixv1.RadixRegistrationSpec{WBS: wbs}}}
		actual := MapContainerBulkDtoFromPod(pod, rrMap, make(map[string]*corev1.LimitRange), clock.NewFakeClock(currentTime))
		assert.Len(t, actual, 1)
		assert.Equal(t, expected, actual)
	})

	// Pod, correct label, container running, rr missing
	t.Run("correct label, container running, rr missing", func(t *testing.T) {
		t.Parallel()
		podName, env, node, app, containerName, containerId, startedAt, currentTime :=
			"pod1", "prod", "node1", "app1", "c1", "cid1", time.Date(2020, 1, 1, 1, 1, 1, 0, time.UTC), time.Date(2020, 2, 1, 1, 1, 1, 0, time.UTC)
		pod := buildPodForTest(podName, fmt.Sprintf("%s-%s", app, env), node,
			setPodAppLabel(app),
			appendPodContainerStatus(
				buildContainerStatusForTest(containerName, containerId,
					setContainerStateRunning(startedAt),
				),
			),
		)
		expected := []repository.ContainerBulkDto{{
			ContainerID:          containerId,
			ContainerName:        containerName,
			PodName:              podName,
			ApplicationName:      app,
			EnvironmentName:      env,
			ComponentName:        "",
			Wbs:                  "",
			StartedAt:            startedAt,
			LastKnowRunningAt:    currentTime,
			CPURequestMillicores: 0,
			MemoryRequestBytes:   0,
			NodeName:             node,
		}}
		actual := MapContainerBulkDtoFromPod(pod, make(map[string]*radixv1.RadixRegistration), make(map[string]*corev1.LimitRange), clock.NewFakeClock(currentTime))
		assert.Len(t, actual, 1)
		assert.Equal(t, expected, actual)
	})

	// Pod, correct label, container terminated
	t.Run("correct label, container terminated", func(t *testing.T) {
		t.Parallel()
		podName, env, node, app, containerName, containerId, startedAt, finishedAt :=
			"pod1", "prod", "node1", "app1", "c1", "cid1", time.Date(2020, 1, 1, 1, 1, 1, 0, time.UTC), time.Date(2020, 2, 1, 1, 1, 1, 0, time.UTC)
		pod := buildPodForTest(podName, fmt.Sprintf("%s-%s", app, env), node,
			setPodAppLabel(app),
			appendPodContainerStatus(
				buildContainerStatusForTest(containerName, containerId,
					setContainerStateTerminated(containerId, startedAt, finishedAt),
				),
			),
		)
		expected := []repository.ContainerBulkDto{{
			ContainerID:          containerId,
			ContainerName:        containerName,
			PodName:              podName,
			ApplicationName:      app,
			EnvironmentName:      env,
			ComponentName:        "",
			Wbs:                  "",
			StartedAt:            startedAt,
			LastKnowRunningAt:    finishedAt,
			CPURequestMillicores: 0,
			MemoryRequestBytes:   0,
			NodeName:             node,
		}}
		actual := MapContainerBulkDtoFromPod(pod, make(map[string]*radixv1.RadixRegistration), make(map[string]*corev1.LimitRange), &clock.RealClock{})
		assert.Len(t, actual, 1)
		assert.Equal(t, expected, actual)
	})

	// Pod, correct label, container terminated
	t.Run("correct label, container terminated, last termination state with different ID", func(t *testing.T) {
		t.Parallel()
		podName, env, node, app, containerName, containerId, startedAt, finishedAt :=
			"pod1", "prod", "node1", "app1", "c1", "cid1", time.Date(2020, 1, 1, 1, 1, 1, 0, time.UTC), time.Date(2020, 2, 1, 1, 1, 1, 0, time.UTC)
		lastContainerId, lastStartedAt, lastFinishedAt := "cid2", time.Date(2021, 1, 1, 1, 1, 1, 0, time.UTC), time.Date(2021, 2, 1, 1, 1, 1, 0, time.UTC)
		pod := buildPodForTest(podName, fmt.Sprintf("%s-%s", app, env), node,
			setPodAppLabel(app),
			appendPodContainerStatus(
				buildContainerStatusForTest(containerName, containerId,
					setContainerStateTerminated(containerId, startedAt, finishedAt),
					setContainerLastTerminatedState(lastContainerId, lastStartedAt, lastFinishedAt),
				),
			),
		)
		expected := []repository.ContainerBulkDto{{
			ContainerID:          containerId,
			ContainerName:        containerName,
			PodName:              podName,
			ApplicationName:      app,
			EnvironmentName:      env,
			ComponentName:        "",
			Wbs:                  "",
			StartedAt:            startedAt,
			LastKnowRunningAt:    finishedAt,
			CPURequestMillicores: 0,
			MemoryRequestBytes:   0,
			NodeName:             node,
		}, {
			ContainerID:          lastContainerId,
			ContainerName:        containerName,
			PodName:              podName,
			ApplicationName:      app,
			EnvironmentName:      env,
			ComponentName:        "",
			Wbs:                  "",
			StartedAt:            lastStartedAt,
			LastKnowRunningAt:    lastFinishedAt,
			CPURequestMillicores: 0,
			MemoryRequestBytes:   0,
			NodeName:             node,
		}}
		actual := MapContainerBulkDtoFromPod(pod, make(map[string]*radixv1.RadixRegistration), make(map[string]*corev1.LimitRange), &clock.RealClock{})
		assert.ElementsMatch(t, expected, actual)
	})

	t.Run("correct label, container terminated, last termination state with same ID", func(t *testing.T) {
		t.Parallel()
		podName, env, node, app, containerName, containerId, startedAt, finishedAt :=
			"pod1", "prod", "node1", "app1", "c1", "cid1", time.Date(2020, 1, 1, 1, 1, 1, 0, time.UTC), time.Date(2020, 2, 1, 1, 1, 1, 0, time.UTC)
		lastContainerId, lastStartedAt, lastFinishedAt := containerId, time.Date(2021, 1, 1, 1, 1, 1, 0, time.UTC), time.Date(2021, 2, 1, 1, 1, 1, 0, time.UTC)
		pod := buildPodForTest(podName, fmt.Sprintf("%s-%s", app, env), node,
			setPodAppLabel(app),
			appendPodContainerStatus(
				buildContainerStatusForTest(containerName, containerId,
					setContainerStateTerminated(containerId, startedAt, finishedAt),
					setContainerLastTerminatedState(lastContainerId, lastStartedAt, lastFinishedAt),
				),
			),
		)
		expected := []repository.ContainerBulkDto{{
			ContainerID:          lastContainerId,
			ContainerName:        containerName,
			PodName:              podName,
			ApplicationName:      app,
			EnvironmentName:      env,
			ComponentName:        "",
			Wbs:                  "",
			StartedAt:            lastStartedAt,
			LastKnowRunningAt:    lastFinishedAt,
			CPURequestMillicores: 0,
			MemoryRequestBytes:   0,
			NodeName:             node,
		}}
		actual := MapContainerBulkDtoFromPod(pod, make(map[string]*radixv1.RadixRegistration), make(map[string]*corev1.LimitRange), &clock.RealClock{})
		assert.ElementsMatch(t, expected, actual)
	})

	// Pod, correct label, container waiting
	t.Run("correct label, container waiting", func(t *testing.T) {
		t.Parallel()
		podName, env, node, app, containerName, containerId :=
			"pod1", "prod", "node1", "app1", "c1", "cid1"
		pod := buildPodForTest(podName, fmt.Sprintf("%s-%s", app, env), node,
			setPodAppLabel(app),
			appendPodContainerStatus(
				buildContainerStatusForTest(containerName, containerId,
					setContainerStateWaiting(),
				),
			),
		)
		actual := MapContainerBulkDtoFromPod(pod, make(map[string]*radixv1.RadixRegistration), make(map[string]*corev1.LimitRange), &clock.RealClock{})
		assert.Len(t, actual, 0)
	})

	// Pod, correct label, container missing Id
	t.Run("correct label, container missing Id", func(t *testing.T) {
		t.Parallel()
		podName, env, node, app, containerName, containerId, startedAt, finishedAt :=
			"pod1", "prod", "node1", "app1", "c1", "", time.Date(2020, 1, 1, 1, 1, 1, 0, time.UTC), time.Date(2020, 2, 1, 1, 1, 1, 0, time.UTC)
		pod := buildPodForTest(podName, fmt.Sprintf("%s-%s", app, env), node,
			setPodAppLabel(app),
			appendPodContainerStatus(
				buildContainerStatusForTest(containerName, containerId,
					setContainerStateTerminated(containerId, startedAt, finishedAt),
				),
			),
		)
		actual := MapContainerBulkDtoFromPod(pod, make(map[string]*radixv1.RadixRegistration), make(map[string]*corev1.LimitRange), &clock.RealClock{})
		assert.Len(t, actual, 0)
	})

	// Pod, correct label, missing resources, lr exist
	t.Run("correct label, missing resources, lr exist", func(t *testing.T) {
		t.Parallel()
		podName, env, node, app, containerName, containerId, millicores, memory, startedAt, currentTime :=
			"pod1", "prod", "node1", "app1", "c1", "cid1", int64(100), int64(1024), time.Date(2020, 1, 1, 1, 1, 1, 0, time.UTC), time.Date(2020, 2, 1, 1, 1, 1, 0, time.UTC)
		pod := buildPodForTest(podName, fmt.Sprintf("%s-%s", app, env), node,
			setPodAppLabel(app),
			appendPodContainer(
				buildContainerForTest(containerName),
			),
			appendPodContainerStatus(
				buildContainerStatusForTest(containerName, containerId,
					setContainerStateRunning(startedAt),
				),
			),
		)
		expected := []repository.ContainerBulkDto{{
			ContainerID:          containerId,
			ContainerName:        containerName,
			PodName:              podName,
			ApplicationName:      app,
			EnvironmentName:      env,
			ComponentName:        "",
			Wbs:                  "",
			StartedAt:            startedAt,
			LastKnowRunningAt:    currentTime,
			CPURequestMillicores: millicores,
			MemoryRequestBytes:   memory,
			NodeName:             node,
		}}
		lrMap := map[string]*corev1.LimitRange{
			fmt.Sprintf("%s-%s", app, env): {
				Spec: corev1.LimitRangeSpec{
					Limits: []corev1.LimitRangeItem{
						{
							Type: corev1.LimitTypeContainer, DefaultRequest: corev1.ResourceList{
								corev1.ResourceCPU:    *resource.NewMilliQuantity(millicores, resource.DecimalSI),
								corev1.ResourceMemory: *resource.NewQuantity(memory, resource.DecimalSI),
							},
						},
					},
				},
			},
		}
		actual := MapContainerBulkDtoFromPod(pod, make(map[string]*radixv1.RadixRegistration), lrMap, clock.NewFakeClock(currentTime))
		assert.Len(t, actual, 1)
		assert.Equal(t, expected, actual)
	})

	// Pod, correct label, missing resources, missing lr
	t.Run("correct label, missing resources, missing lr", func(t *testing.T) {
		t.Parallel()
		podName, env, node, app, containerName, containerId, startedAt, currentTime :=
			"pod1", "prod", "node1", "app1", "c1", "cid1", time.Date(2020, 1, 1, 1, 1, 1, 0, time.UTC), time.Date(2020, 2, 1, 1, 1, 1, 0, time.UTC)
		pod := buildPodForTest(podName, fmt.Sprintf("%s-%s", app, env), node,
			setPodAppLabel(app),
			appendPodContainer(
				buildContainerForTest(containerName),
			),
			appendPodContainerStatus(
				buildContainerStatusForTest(containerName, containerId,
					setContainerStateRunning(startedAt),
				),
			),
		)
		expected := []repository.ContainerBulkDto{{
			ContainerID:          containerId,
			ContainerName:        containerName,
			PodName:              podName,
			ApplicationName:      app,
			EnvironmentName:      env,
			ComponentName:        "",
			Wbs:                  "",
			StartedAt:            startedAt,
			LastKnowRunningAt:    currentTime,
			CPURequestMillicores: 0,
			MemoryRequestBytes:   0,
			NodeName:             node,
		}}

		actual := MapContainerBulkDtoFromPod(pod, make(map[string]*radixv1.RadixRegistration), make(map[string]*corev1.LimitRange), clock.NewFakeClock(currentTime))
		assert.Len(t, actual, 1)
		assert.Equal(t, expected, actual)
	})

	// Pod, correct label, cpu resources, lr memory
	t.Run("correct label, cpu resources, lr memory", func(t *testing.T) {
		t.Parallel()
		podName, env, node, app, containerName, containerId, millicores, memory, startedAt, currentTime :=
			"pod1", "prod", "node1", "app1", "c1", "cid1", int64(100), int64(1024), time.Date(2020, 1, 1, 1, 1, 1, 0, time.UTC), time.Date(2020, 2, 1, 1, 1, 1, 0, time.UTC)
		pod := buildPodForTest(podName, fmt.Sprintf("%s-%s", app, env), node,
			setPodAppLabel(app),
			appendPodContainer(
				buildContainerForTest(containerName,
					setContainerResourceRequest(corev1.ResourceCPU, *resource.NewMilliQuantity(millicores, resource.DecimalSI)),
				),
			),
			appendPodContainerStatus(
				buildContainerStatusForTest(containerName, containerId,
					setContainerStateRunning(startedAt),
				),
			),
		)
		expected := []repository.ContainerBulkDto{{
			ContainerID:          containerId,
			ContainerName:        containerName,
			PodName:              podName,
			ApplicationName:      app,
			EnvironmentName:      env,
			ComponentName:        "",
			Wbs:                  "",
			StartedAt:            startedAt,
			LastKnowRunningAt:    currentTime,
			CPURequestMillicores: millicores,
			MemoryRequestBytes:   memory,
			NodeName:             node,
		}}
		lrMap := map[string]*corev1.LimitRange{
			fmt.Sprintf("%s-%s", app, env): {
				Spec: corev1.LimitRangeSpec{
					Limits: []corev1.LimitRangeItem{
						{
							Type: corev1.LimitTypeContainer, DefaultRequest: corev1.ResourceList{
								corev1.ResourceMemory: *resource.NewQuantity(memory, resource.DecimalSI),
							},
						},
					},
				},
			},
		}
		actual := MapContainerBulkDtoFromPod(pod, make(map[string]*radixv1.RadixRegistration), lrMap, clock.NewFakeClock(currentTime))
		assert.Len(t, actual, 1)
		assert.Equal(t, expected, actual)
	})

	// Pod, correct label, memory resources, lr cpu
	t.Run("correct label, memory resources, lr cpu", func(t *testing.T) {
		t.Parallel()
		podName, env, node, app, containerName, containerId, millicores, memory, startedAt, currentTime :=
			"pod1", "prod", "node1", "app1", "c1", "cid1", int64(100), int64(1024), time.Date(2020, 1, 1, 1, 1, 1, 0, time.UTC), time.Date(2020, 2, 1, 1, 1, 1, 0, time.UTC)
		pod := buildPodForTest(podName, fmt.Sprintf("%s-%s", app, env), node,
			setPodAppLabel(app),
			appendPodContainer(
				buildContainerForTest(containerName,
					setContainerResourceRequest(corev1.ResourceMemory, *resource.NewQuantity(memory, resource.DecimalSI)),
				),
			),
			appendPodContainerStatus(
				buildContainerStatusForTest(containerName, containerId,
					setContainerStateRunning(startedAt),
				),
			),
		)
		expected := []repository.ContainerBulkDto{{
			ContainerID:          containerId,
			ContainerName:        containerName,
			PodName:              podName,
			ApplicationName:      app,
			EnvironmentName:      env,
			ComponentName:        "",
			Wbs:                  "",
			StartedAt:            startedAt,
			LastKnowRunningAt:    currentTime,
			CPURequestMillicores: millicores,
			MemoryRequestBytes:   memory,
			NodeName:             node,
		}}
		lrMap := map[string]*corev1.LimitRange{
			fmt.Sprintf("%s-%s", app, env): {
				Spec: corev1.LimitRangeSpec{
					Limits: []corev1.LimitRangeItem{
						{
							Type: corev1.LimitTypeContainer, DefaultRequest: corev1.ResourceList{
								corev1.ResourceCPU: *resource.NewMilliQuantity(millicores, resource.DecimalSI),
							},
						},
					},
				},
			},
		}
		actual := MapContainerBulkDtoFromPod(pod, make(map[string]*radixv1.RadixRegistration), lrMap, clock.NewFakeClock(currentTime))
		assert.Len(t, actual, 1)
		assert.Equal(t, expected, actual)
	})
	// Pod, correct label, two containers
	t.Run("correct label, two containers", func(t *testing.T) {
		t.Parallel()
		podName, env, node, app, containerName, containerId1, containerId2, startedAt, currentTime :=
			"pod1", "prod", "node1", "app1", "c1", "cid1", "cid1", time.Date(2020, 1, 1, 1, 1, 1, 0, time.UTC), time.Date(2020, 2, 1, 1, 1, 1, 0, time.UTC)
		pod := buildPodForTest(podName, fmt.Sprintf("%s-%s", app, env), node,
			setPodAppLabel(app),
			appendPodContainerStatus(
				buildContainerStatusForTest(containerName, containerId1,
					setContainerStateRunning(startedAt),
				),
			),
			appendPodContainerStatus(
				buildContainerStatusForTest(containerName, containerId2,
					setContainerStateRunning(startedAt),
				),
			),
		)
		expected := []repository.ContainerBulkDto{
			{
				ContainerID:          containerId1,
				ContainerName:        containerName,
				PodName:              podName,
				ApplicationName:      app,
				EnvironmentName:      env,
				ComponentName:        "",
				Wbs:                  "",
				StartedAt:            startedAt,
				LastKnowRunningAt:    currentTime,
				CPURequestMillicores: 0,
				MemoryRequestBytes:   0,
				NodeName:             node,
			},
			{
				ContainerID:          containerId2,
				ContainerName:        containerName,
				PodName:              podName,
				ApplicationName:      app,
				EnvironmentName:      env,
				ComponentName:        "",
				Wbs:                  "",
				StartedAt:            startedAt,
				LastKnowRunningAt:    currentTime,
				CPURequestMillicores: 0,
				MemoryRequestBytes:   0,
				NodeName:             node,
			},
		}
		actual := MapContainerBulkDtoFromPod(pod, make(map[string]*radixv1.RadixRegistration), make(map[string]*corev1.LimitRange), clock.NewFakeClock(currentTime))
		assert.Len(t, actual, 2)
		assert.ElementsMatch(t, expected, actual)
	})

	// Pod, correct label, last termination state set
	t.Run("correct label, missing resources, missing lr", func(t *testing.T) {
		t.Parallel()
		podName, env, node, app, containerName, containerId, terminatedContainerId, startedAt, currentTime, terminatedStartedAt, teminatedFinishedAt :=
			"pod1", "prod", "node1", "app1", "c1", "cid1", "cid2",
			time.Date(2020, 1, 1, 1, 1, 1, 0, time.UTC), time.Date(2020, 2, 1, 1, 1, 1, 0, time.UTC),
			time.Date(2020, 3, 1, 1, 1, 1, 0, time.UTC), time.Date(2020, 4, 1, 1, 1, 1, 0, time.UTC)
		pod := buildPodForTest(podName, fmt.Sprintf("%s-%s", app, env), node,
			setPodAppLabel(app),
			appendPodContainer(
				buildContainerForTest(containerName),
			),
			appendPodContainerStatus(
				buildContainerStatusForTest(containerName, containerId,
					setContainerStateRunning(startedAt),
					setContainerLastTerminatedState(terminatedContainerId, terminatedStartedAt, teminatedFinishedAt),
				),
			),
		)
		expected := []repository.ContainerBulkDto{
			{
				ContainerID:          containerId,
				ContainerName:        containerName,
				PodName:              podName,
				ApplicationName:      app,
				EnvironmentName:      env,
				ComponentName:        "",
				Wbs:                  "",
				StartedAt:            startedAt,
				LastKnowRunningAt:    currentTime,
				CPURequestMillicores: 0,
				MemoryRequestBytes:   0,
				NodeName:             node,
			},
			{
				ContainerID:          terminatedContainerId,
				ContainerName:        containerName,
				PodName:              podName,
				ApplicationName:      app,
				EnvironmentName:      env,
				ComponentName:        "",
				Wbs:                  "",
				StartedAt:            terminatedStartedAt,
				LastKnowRunningAt:    teminatedFinishedAt,
				CPURequestMillicores: 0,
				MemoryRequestBytes:   0,
				NodeName:             node,
			},
		}

		actual := MapContainerBulkDtoFromPod(pod, make(map[string]*radixv1.RadixRegistration), make(map[string]*corev1.LimitRange), clock.NewFakeClock(currentTime))
		assert.Len(t, actual, 2)
		assert.ElementsMatch(t, expected, actual)
	})

	// Pod, incorrect label
	t.Run("incorrect label", func(t *testing.T) {
		t.Parallel()
		podName, env, node, app, containerName, containerId, startedAt, currentTime :=
			"pod1", "prod", "node1", "app1", "c1", "cid1", time.Date(2020, 1, 1, 1, 1, 1, 0, time.UTC), time.Date(2020, 2, 1, 1, 1, 1, 0, time.UTC)
		pod := buildPodForTest(podName, fmt.Sprintf("%s-%s", app, env), node,
			appendPodContainerStatus(
				buildContainerStatusForTest(containerName, containerId,
					setContainerStateRunning(startedAt),
				),
			),
		)

		actual := MapContainerBulkDtoFromPod(pod, make(map[string]*radixv1.RadixRegistration), make(map[string]*corev1.LimitRange), clock.NewFakeClock(currentTime))
		assert.Len(t, actual, 0)
	})
}

func setPodAppLabel(appName string) func(*corev1.Pod) {
	return func(p *corev1.Pod) {
		p.Labels[kube.RadixAppLabel] = appName
	}
}

func setPodComponentLabel(componentName string) func(*corev1.Pod) {
	return func(p *corev1.Pod) {
		p.Labels[kube.RadixComponentLabel] = componentName
	}
}

func appendPodContainerStatus(s corev1.ContainerStatus) func(*corev1.Pod) {
	return func(p *corev1.Pod) {
		p.Status.ContainerStatuses = append(p.Status.ContainerStatuses, s)
	}
}

func appendPodContainer(c corev1.Container) func(*corev1.Pod) {
	return func(p *corev1.Pod) {
		p.Spec.Containers = append(p.Spec.Containers, c)
	}
}

func setContainerStateWaiting() func(*corev1.ContainerStatus) {
	return func(cs *corev1.ContainerStatus) {
		cs.State = corev1.ContainerState{Waiting: &corev1.ContainerStateWaiting{}}
	}
}

func setContainerStateRunning(startedAt time.Time) func(*corev1.ContainerStatus) {
	return func(cs *corev1.ContainerStatus) {
		cs.State = corev1.ContainerState{Running: &corev1.ContainerStateRunning{StartedAt: metav1.NewTime(startedAt)}}
	}
}

func setContainerStateTerminated(containerId string, startedAt time.Time, finishedAt time.Time) func(*corev1.ContainerStatus) {
	return func(cs *corev1.ContainerStatus) {
		cs.State = corev1.ContainerState{Terminated: &corev1.ContainerStateTerminated{
			ContainerID: containerId,
			StartedAt:   metav1.NewTime(startedAt),
			FinishedAt:  metav1.NewTime(finishedAt),
		}}
	}
}

func setContainerLastTerminatedState(containerId string, startedAt time.Time, finishedAt time.Time) func(*corev1.ContainerStatus) {
	return func(cs *corev1.ContainerStatus) {
		cs.LastTerminationState = corev1.ContainerState{Terminated: &corev1.ContainerStateTerminated{
			ContainerID: containerId,
			StartedAt:   metav1.NewTime(startedAt),
			FinishedAt:  metav1.NewTime(finishedAt),
		}}
	}
}

func buildContainerStatusForTest(name, id string, builders ...func(*corev1.ContainerStatus)) corev1.ContainerStatus {
	cs := corev1.ContainerStatus{Name: name, ContainerID: id}

	for _, b := range builders {
		b(&cs)
	}

	return cs
}

func setContainerResourceRequest(resource corev1.ResourceName, qty resource.Quantity) func(*corev1.Container) {
	return func(c *corev1.Container) {
		if c.Resources.Requests == nil {
			c.Resources.Requests = make(corev1.ResourceList)
		}
		c.Resources.Requests[resource] = qty
	}
}

func buildContainerForTest(name string, builders ...func(*corev1.Container)) corev1.Container {
	c := corev1.Container{Name: name}

	for _, b := range builders {
		b(&c)
	}

	return c
}

func buildPodForTest(name, namespace, node string, builders ...func(*corev1.Pod)) *corev1.Pod {
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace, Labels: make(map[string]string)},
		Spec: corev1.PodSpec{
			NodeName: node,
		},
	}

	for _, b := range builders {
		b(pod)
	}

	return pod
}
