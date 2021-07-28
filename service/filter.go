package service

import (
	"k8s.io/api/core/v1"
)

const (
	CpuPriorityPred = "cpuPriorityPred"
)

type FitPredicate func(pod *v1.Pod, node v1.Node) (bool, []string, error)

var predicatesFuncs = map[string]FitPredicate{
	CpuPriorityPred: cpuPriorityPredicate,
}

//podFitsOneNode Filter out all nodes suitable for pod
func podFitsOneNode(pod *v1.Pod, node v1.Node) (bool, []string, error) {
	return predicatesFuncs[CpuPriorityPred](pod, node)
}

func cpuPriorityPredicate(pod *v1.Pod, node v1.Node) (bool, []string, error) {
	failReasons := make([]string, 1)
	cpuSize := 0
	for _, c := range pod.Spec.Containers {
		cpuSize += c.Resources.Limits.Cpu().Size()
	}

	if cpuSize > 10000 && *pod.Spec.Priority < 10000000 {
		failReasons = append(failReasons, "Cannot be scheduled, pod size is greater than 10000, and priority is less than 10000000")
		return false, failReasons, nil
	}
	return true, failReasons, nil
}
