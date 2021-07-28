package router

import (
	"github.com/emicklei/go-restful"
	"fmt"
	schedulerapi "k8s.io/kube-scheduler/extender/v1"
	"encoding/json"
	"bytes"
	"io"
	"github.com/golang/glog"
	"github.com/custom-scheduler-extender/service"
	"net/http"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func Index(r *restful.Request, w *restful.Response) {
	fmt.Fprintln(w.ResponseWriter, "a extender for kube-scheduler!")
}

func predicate(r *restful.Request, w *restful.Response) {
	var buf bytes.Buffer
	var extendArgs schedulerapi.ExtenderArgs
	var extenderFilterResult schedulerapi.ExtenderFilterResult

	body := io.TeeReader(r.Request.Body, &buf)
	if err := json.NewDecoder(body).Decode(&extendArgs); err != nil {
		glog.Error("json decode err:%v", err)
		extenderFilterResult.Error = err.Error()
	} else {
		extenderFilterResult = service.Filter(extendArgs)
	}
	if response, err := json.Marshal(&extenderFilterResult); err != nil {
		glog.Fatal(err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func priority(r *restful.Request, w *restful.Response) {
	var buf bytes.Buffer
	var extendArgs schedulerapi.ExtenderArgs
	var hostPriorityList *schedulerapi.HostPriorityList
	body := io.TeeReader(r.Request.Body, &buf)
	if err := json.NewDecoder(body).Decode(&extendArgs); err != nil {
		glog.Error("json decode err:%v", err)
		hostPriorityList = &schedulerapi.HostPriorityList{}
	} else {
		hostPriorityList = service.Prioritize(extendArgs)
	}
	if response, err := json.Marshal(&hostPriorityList); err != nil {
		glog.Fatal(err)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func bind(r *restful.Request, w *restful.Response) {
	var extenderBindingArgs schedulerapi.ExtenderBindingArgs
	var buf bytes.Buffer
	var extenderBindingResult schedulerapi.ExtenderBindingResult

	body := io.TeeReader(r.Request.Body, &buf)
	if err := json.NewDecoder(body).Decode(&extenderBindingArgs); err != nil {
		extenderBindingResult.Error = err.Error()
	} else {
		extenderBindingResult = service.Bind(v1.Binding{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: extenderBindingArgs.PodNamespace,
				Name:      extenderBindingArgs.PodName,
				UID:       extenderBindingArgs.PodUID,
			},
			Target: v1.ObjectReference{Kind: "Node", Name: extenderBindingArgs.Node},
		})
	}

}
