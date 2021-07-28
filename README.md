## 自定义调度器

## 1. 通过[scheduler_extender](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/scheduling/scheduler_extender.md)的方式进行扩展

在调度 Pod 时，扩展器允许外部进程过滤节点并确定其优先级。 向扩展程序发出两个单独的 http/https 调用，一个用于“过滤器”，另一个用于“优先”操作。 此外，扩展程序可以选择通过实现“绑定”操作将 pod 绑定到 apiserver。

 要使用扩展程序，您必须创建调度程序策略配置文件。 配置指定如何到达扩展器，是使用 http 还是 https 以及超时。 

**即实现两个http接口，一个用于”预选“，一个用于“优选”**



extender配置数据结构如下

```
// Holds the parameters used to communicate with the extender. If a verb is unspecified/empty,
// it is assumed that the extender chose not to provide that extension.
type ExtenderConfig struct {
	// URLPrefix at which the extender is available
	URLPrefix string `json:"urlPrefix"`
	// Verb for the filter call, empty if not supported. This verb is appended to the URLPrefix when issuing the filter call to extender.
	FilterVerb string `json:"filterVerb,omitempty"`
	// Verb for the prioritize call, empty if not supported. This verb is appended to the URLPrefix when issuing the prioritize call to extender.
	PrioritizeVerb string `json:"prioritizeVerb,omitempty"`
	// Verb for the bind call, empty if not supported. This verb is appended to the URLPrefix when issuing the bind call to extender.
	// If this method is implemented by the extender, it is the extender's responsibility to bind the pod to apiserver.
	BindVerb string `json:"bindVerb,omitempty"`
	// The numeric multiplier for the node scores that the prioritize call generates.
	// The weight should be a positive integer
	Weight int `json:"weight,omitempty"`
	// EnableHttps specifies whether https should be used to communicate with the extender
	EnableHttps bool `json:"enableHttps,omitempty"`
	// TLSConfig specifies the transport layer security config
	TLSConfig *client.TLSClientConfig `json:"tlsConfig,omitempty"`
	// HTTPTimeout specifies the timeout duration for a call to the extender. Filter timeout fails the scheduling of the pod. Prioritize
	// timeout is ignored, k8s/other extenders priorities are used to select the node.
	HTTPTimeout time.Duration `json:"httpTimeout,omitempty"`
}
```



配置例子：

```
{
  "predicates": [
    {
      "name": "HostName"
    },
    {
      "name": "MatchNodeSelector"
    },
    {
      "name": "PodFitsResources"
    }
  ],
  "priorities": [
    {
      "name": "LeastRequestedPriority",
      "weight": 1
    }
  ],
  "extenders": [
    {
      "urlPrefix": "http://127.0.0.1:12345/api/scheduler",
      "filterVerb": "filter",
      "enableHttps": false
    }
  ]
}
```



http/https接口入参结构体

```
// ExtenderArgs represents the arguments needed by the extender to filter/prioritize
// nodes for a pod.
type ExtenderArgs struct {
	// Pod being scheduled
	Pod   api.Pod      `json:"pod"`
	// List of candidate nodes where the pod can be scheduled
	Nodes api.NodeList `json:"nodes"`
}
```



“filter”接口返回节点的数组（schedulerapi.ExtenderFilterResult），“prioritize”接口返回每个节点的优先级(schedulerapi.HostPriorityList)，“filter”接口返回的节点数组基于预选阶段，“prioritize”接口返回的优先级数组基于优选打分阶段，作为最终的节点选择。

“bind”接口用于将pod绑定到目标节点，它也可以通过extender来实现。当它实现时。extender向apiserver发出绑定调用。

“bind”接口的入参结构体

```
// ExtenderBindingArgs represents the arguments to an extender for binding a pod to a node.
type ExtenderBindingArgs struct {
	// PodName is the name of the pod being bound
	PodName string
	// PodNamespace is the namespace of the pod being bound
	PodNamespace string
	// PodUID is the UID of the pod being bound
	PodUID types.UID
	// Node selected by the scheduler
	Node string
}
```

