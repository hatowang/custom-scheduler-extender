package main
import (
	restful "github.com/emicklei/go-restful/v3"
	"github.com/custom-scheduler-extender/router"
	"github.com/golang/glog"
	"net/http"
)


func main() {
	restful.Add(router.NewRouters())
	glog.Fatal("start http!", http.ListenAndServe(":8081",nil))
}
