package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/equinor/radix-operator/pkg/apis/kube"
	v1 "github.com/equinor/radix-operator/pkg/apis/radix/v1"
	corev1 "k8s.io/api/core/v1"
)

type ContainerBulkTvp struct {
	ContainerId          string
	ContainerName        string
	PodName              string
	ApplicationName      string
	EnvironmentName      string
	ComponentName        string
	Wbs                  string
	StartedAt            time.Time
	LastKnowRunningAt    time.Time
	CpuRequestMillicores int64
	MemoryRequestBytes   int64
	NodeName             string
}

func ContainerBulkTvpFromPod(pod *corev1.Pod, rrMap map[string]*v1.RadixRegistration, limitRangeMap map[string]*corev1.LimitRange) (containersTvp []ContainerBulkTvp, err error) {
	if pod == nil {
		err = errors.New("pod is nil")
		return
	}

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

		containerTvp := ContainerBulkTvp{
			ContainerId:     containerStatus.ContainerID,
			ContainerName:   containerStatus.Name,
			PodName:         pod.Name,
			ApplicationName: appName,
			EnvironmentName: environmentName,
			ComponentName:   componentName,
			NodeName:        pod.Spec.NodeName,
		}

		if container := getContainerByName(containerStatus.Name, pod.Spec.Containers); container != nil {
			setContainerBulkTvpResourceProps(&containerTvp, container)
		}
		setContainerBulkTvpLimitRangeProps(&containerTvp, limitRangeMap[pod.Namespace])
		setContainerBulkTvpRunningProps(&containerTvp, containerStatus.State.Running)
		setContainerBulkTvpTerminatedProps(&containerTvp, containerStatus.State.Terminated)
		setContainerBulkTvpRadixRegistrationProps(&containerTvp, rrMap[appName])

		containersTvp = append(containersTvp, containerTvp)

		if lastTerminatedState := containerStatus.LastTerminationState.Terminated; lastTerminatedState != nil {
			lastTerminatedTvp := containerTvp
			setContainerBulkTvpTerminatedProps(&lastTerminatedTvp, lastTerminatedState)
			containersTvp = append(containersTvp, lastTerminatedTvp)
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

func setContainerBulkTvpRadixRegistrationProps(cbt *ContainerBulkTvp, rr *v1.RadixRegistration) {
	if cbt == nil || rr == nil {
		return
	}

	cbt.Wbs = rr.Spec.WBS
}

func setContainerBulkTvpResourceProps(cbt *ContainerBulkTvp, container *corev1.Container) {
	if cbt == nil || container == nil {
		return
	}

	setContainerBulkTvpMemory(cbt, container.Resources.Requests)
	setContainerBulkTvpCpu(cbt, container.Resources.Requests)
}

func setContainerBulkTvpCpu(cbt *ContainerBulkTvp, resourceList corev1.ResourceList) {
	if cpu := resourceList.Cpu(); cpu != nil {
		cbt.CpuRequestMillicores = cpu.MilliValue()
	}
}

func setContainerBulkTvpMemory(cbt *ContainerBulkTvp, resourceList corev1.ResourceList) {
	if mem := resourceList.Memory(); mem != nil {
		cbt.MemoryRequestBytes = mem.Value()
	}
}

func setContainerBulkTvpLimitRangeProps(cbt *ContainerBulkTvp, limitRange *corev1.LimitRange) {
	if cbt == nil || limitRange == nil {
		return
	}

	if lri := getFirstContainerLimitRangeItem(limitRange.Spec.Limits); lri != nil {
		if cbt.MemoryRequestBytes == 0 {
			setContainerBulkTvpMemory(cbt, lri.DefaultRequest)
		}

		if cbt.CpuRequestMillicores == 0 {
			setContainerBulkTvpCpu(cbt, lri.DefaultRequest)
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

func setContainerBulkTvpTerminatedProps(cbt *ContainerBulkTvp, terminated *corev1.ContainerStateTerminated) {
	if cbt == nil || terminated == nil {
		return
	}

	cbt.ContainerId = terminated.ContainerID
	cbt.StartedAt = terminated.StartedAt.Time
	cbt.LastKnowRunningAt = terminated.FinishedAt.Time
}

func setContainerBulkTvpRunningProps(cbt *ContainerBulkTvp, running *corev1.ContainerStateRunning) {
	if cbt == nil || running == nil {
		return
	}

	cbt.StartedAt = running.StartedAt.Time
	cbt.LastKnowRunningAt = time.Now()
}

func getContainerByName(name string, containers []corev1.Container) *corev1.Container {
	for _, container := range containers {
		if container.Name == name {
			return &container
		}
	}

	return nil
}
