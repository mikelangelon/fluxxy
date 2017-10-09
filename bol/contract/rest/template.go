package rest

import (
	"bol/contract/contract_template"
	"github.com/emicklei/go-restful"
)

var (
	TemplateWebService = &restful.WebService{}
)

func init() {
	id := TemplateWebService.PathParameter("id", "Identifier of the resource").DataType("long")

	TemplateWebService.
		Path("/v1/contract_template").
		Doc("Template related").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML).
		ApiVersion("1.0")

	TemplateWebService.Route(
		TemplateWebService.POST("/").
			Doc("Creates a new Template").
			Reads(contract_template.Template{}).
			To(createTemplate))

	TemplateWebService.Route(
		TemplateWebService.PUT("/{id}").
			Param(id).
			Doc("Update a new Template").
			Reads(contract_template.Template{}).
			To(updateTemplate))

	TemplateWebService.Route(
		TemplateWebService.GET("/{id}").
			Param(id).
			Doc("Retrieves a template given an id").
			Writes(contract_template.Template{}).
			To(retrieveTemplate))

}
