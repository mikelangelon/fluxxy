package rest

import (
	"github.com/emicklei/go-restful"
)

var (
	PdfWebService = &restful.WebService{}
)

func init() {
	templateId := TemplateWebService.PathParameter("id", "Identifier of the resource").DataType("long")

	PdfWebService.
		Path("/v1/pdf").
		Doc("Generate pdf").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML).
		ApiVersion("1.0")

	PdfWebService.Route(
		PdfWebService.POST("/generatePdf").
			Doc("Test to generate dummy pdf").
			Produces("application/pdf").
			To(generatePDF))

	PdfWebService.Route(
		PdfWebService.POST("/generatePdf/{id}").
			Doc("Generates the pdf of a contract given its Id").
			Param(templateId).
			Produces("application/pdf").
			To(generatePDFFromTemplate))

}
