package rest

import (
	"github.com/emicklei/go-restful"
	"bol/contract/contract"
)

var (
	ContractWebService = &restful.WebService{}
)

func init() {
	id := TemplateWebService.PathParameter("id", "Identifier of the resource").DataType("long")

	ContractWebService.
	Path("/v1/contract").
		Doc("Contract related").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML).
		ApiVersion("1.0")

	ContractWebService.Route(
		ContractWebService.POST("/").
			Doc("Creates a new contract").
			Reads(contract.Contract{}).
			To(createContract))
	ContractWebService.Route(
		ContractWebService.PUT("/{id}").
			Param(id).
			Doc("Update a new contract").
			Reads(contract.Contract{}).
			To(updateContract))
	ContractWebService.Route(
		ContractWebService.GET("/{id}").
			Param(id).
			Doc("Retrieves a contract given an id").
			Writes(contract.Contract{}).
			To(retrieveContract))
}
