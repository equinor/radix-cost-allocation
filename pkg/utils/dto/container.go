package dto

import (
	"fmt"
	"strings"

	"github.com/equinor/radix-cost-allocation/pkg/repository"
	"github.com/equinor/radix-cost-allocation/pkg/utils/clock"
	"github.com/equinor/radix-operator/pkg/apis/kube"
	radixv1 "github.com/equinor/radix-operator/pkg/apis/radix/v1"
	corev1 "k8s.io/api/core/v1"
)

// MapContainerBulkDtoFromPod builds a ContainerBulkDto from containers in the pod.
// Container information is only extracted if the pod has a "radix-app" label.
// WBS is extracted from the rrMap, where rrMap key must match the value of the "radix-app" label for the pod
// CPU and memory is read from limitRangeMap, where key must match the namespace of the pod, if missing in pod container spec.
func MapContainerBulkDtoFromPod(pod *corev1.Pod, rrMap map[string]*radixv1.RadixRegistration, limitRangeMap map[string]*corev1.LimitRange, clock clock.Clock) (containersDto []repository.ContainerBulkDto) {
	appName, ok := pod.Labels[kube.RadixAppLabel]
	if !ok {
		return
	}

	componentName := pod.Labels[kube.RadixComponentLabel]
	environmentName := getEnvironmentNameFromNamespace(appName, pod.Namespace)

	for _, containerStatus := range pod.Status.ContainerStatuses {
		if containerStatus.State.Waiting != nil {
			continue
		}

		containerDto := repository.ContainerBulkDto{
			ContainerID:     containerStatus.ContainerID,
			ContainerName:   containerStatus.Name,
			PodName:         pod.Name,
			ApplicationName: appName,
			EnvironmentName: environmentName,
			ComponentName:   componentName,
			NodeName:        pod.Spec.NodeName,
		}

		if container := getContainerByName(containerStatus.Name, pod.Spec.Containers); container != nil {
			setContainerBulkDtoResourceProps(&containerDto, container)
		}
		setContainerBulkDtoLimitRangeProps(&containerDto, limitRangeMap[pod.Namespace])
		setContainerBulkDtoRunningProps(&containerDto, containerStatus.State.Running, clock)
		setContainerBulkDtoTerminatedProps(&containerDto, containerStatus.State.Terminated)
		setContainerBulkDtoRadixRegistrationProps(&containerDto, rrMap[appName])

		containersDto = append(containersDto, containerDto)

		if lastTerminatedState := containerStatus.LastTerminationState.Terminated; lastTerminatedState != nil {
			lastTerminatedDto := containerDto
			setContainerBulkDtoTerminatedProps(&lastTerminatedDto, lastTerminatedState)
			containersDto = append(containersDto, lastTerminatedDto)
		}
	}

	return
}

func getEnvironmentNameFromNamespace(appName, ns string) string {
	if env := strings.TrimPrefix(ns, fmt.Sprintf("%s-", appName)); env != ns {
		return env
	}

	return ""
}

func setContainerBulkDtoRadixRegistrationProps(cbt *repository.ContainerBulkDto, rr *radixv1.RadixRegistration) {
	if cbt == nil || rr == nil {
		return
	}

	cbt.Wbs = rr.Spec.WBS
}

func setContainerBulkDtoResourceProps(cbt *repository.ContainerBulkDto, container *corev1.Container) {
	if cbt == nil || container == nil {
		return
	}

	setContainerBulkDtoMemory(cbt, container.Resources.Requests)
	setContainerBulkDtoCPU(cbt, container.Resources.Requests)
}

func setContainerBulkDtoCPU(cbt *repository.ContainerBulkDto, resourceList corev1.ResourceList) {
	if cpu := resourceList.Cpu(); cpu != nil {
		cbt.CPURequestMillicores = cpu.MilliValue()
	}
}

func setContainerBulkDtoMemory(cbt *repository.ContainerBulkDto, resourceList corev1.ResourceList) {
	if mem := resourceList.Memory(); mem != nil {
		cbt.MemoryRequestBytes = mem.Value()
	}
}

func setContainerBulkDtoLimitRangeProps(cbt *repository.ContainerBulkDto, limitRange *corev1.LimitRange) {
	if cbt == nil || limitRange == nil {
		return
	}

	if lri := getFirstContainerLimitRangeItem(limitRange.Spec.Limits); lri != nil {
		if cbt.MemoryRequestBytes == 0 {
			setContainerBulkDtoMemory(cbt, lri.DefaultRequest)
		}

		if cbt.CPURequestMillicores == 0 {
			setContainerBulkDtoCPU(cbt, lri.DefaultRequest)
		}
	}
}

func getFirstContainerLimitRangeItem(items []corev1.LimitRangeItem) *corev1.LimitRangeItem {
	for _, lri := range items {
		if lri.Type == corev1.LimitTypeContainer {
			return &lri
		}
	}

	return nil
}

func setContainerBulkDtoTerminatedProps(cbt *repository.ContainerBulkDto, terminated *corev1.ContainerStateTerminated) {
	if cbt == nil || terminated == nil {
		return
	}

	cbt.ContainerID = terminated.ContainerID
	cbt.StartedAt = terminated.StartedAt.Time
	cbt.LastKnowRunningAt = terminated.FinishedAt.Time
}

func setContainerBulkDtoRunningProps(cbt *repository.ContainerBulkDto, running *corev1.ContainerStateRunning, clock clock.Clock) {
	if cbt == nil || running == nil {
		return
	}

	cbt.StartedAt = running.StartedAt.Time
	cbt.LastKnowRunningAt = clock.Now()
}

func getContainerByName(name string, containers []corev1.Container) *corev1.Container {
	for _, container := range containers {
		if container.Name == name {
			return &container
		}
	}

	return nil
}
