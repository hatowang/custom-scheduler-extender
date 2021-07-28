package router

import (
	restful "github.com/emicklei/go-restful/v3"
)

func NewRouters() *restful.WebService {
	service := new(restful.WebService)
	service.
		Path("/").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_XML, restful.MIME_JSON)

	service.Route(service.GET("/").To(Index))
	service.Route(service.POST("/filter").To(predicate))
	service.Route(service.POST("/prioritize").To(priority))
	service.Route(service.POST("/bind").To(bind))

	return service
}