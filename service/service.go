package service

import (
	schedulerapi "k8s.io/kube-scheduler/extender/v1"
	"k8s.io/api/core/v1"
	"strings"
	"github.com/golang/glog"
)

func Filter(args schedulerapi.ExtenderArgs) schedulerapi.ExtenderFilterResult {
	var filteredNodes []v1.Node
	failedNodes := make(schedulerapi.FailedNodesMap)
	var errStr []string
	for _, node := range args.Nodes.Items {
		fits, failReasons, err := podFitsOneNode(args.Pod, node)
		if fits {
			filteredNodes = append(filteredNodes, node)
		} else {
			failedNodes[node.Name] = strings.Join(failReasons, ",")
		}

		if err != nil {
			errStr = append(errStr, err.Error())
		}
	}

	return schedulerapi.ExtenderFilterResult{
		Nodes: &v1.NodeList{
			Items: filteredNodes,
		},
		FailedNodes: failedNodes,
		Error:       strings.Join(errStr, ","),
	}
}

func Prioritize(args schedulerapi.ExtenderArgs) *schedulerapi.HostPriorityList {
	nodes := args.Nodes.Items
	hostPriorityList := make(schedulerapi.HostPriorityList, len(nodes))
	var score int
	for i, node := range nodes {
		score = i % 10
		hostPriorityList[i] = schedulerapi.HostPriority{
			Host:  node.Name,
			Score: int64(score),
		}
	}
	return &hostPriorityList
}

func Bind(binding v1.Binding) schedulerapi.ExtenderBindingResult {
	//Todo bind
	glog.Infoln("pod: %s, nod:%s, bind success!", binding.Namespace+"/"+binding.Name, binding.Target.Name)
	return schedulerapi.ExtenderBindingResult{
		Error: "",
	}
}
