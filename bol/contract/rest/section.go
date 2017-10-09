package rest

import (
	"github.com/emicklei/go-restful"
	"bol/contract/section"
)

var (
	SectionWebService = &restful.WebService{}
)

func init() {

	id := TemplateWebService.PathParameter("id", "Identifier of the resource").DataType("long")

	SectionWebService.
	Path("/v1/section").
		Doc("Section related").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML).
		ApiVersion("1.0")

	SectionWebService.Route(
		SectionWebService.POST("/").
			Doc("Creates a new section").
			Reads(section.Section{}).
			To(createSection))

	SectionWebService.Route(
		SectionWebService.GET("/{id}").Param(id).
			Doc("Retrieves a section").
			Writes(section.Section{}).
			To(retrieveSection))

	SectionWebService.Route(
		SectionWebService.PUT("/{id}").Param(id).
			Doc("Update a section").
			Reads(section.Section{}).
			To(updateSection))

}
