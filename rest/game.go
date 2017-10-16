package rest

import (
	"github.com/emicklei/go-restful"
)

var (
	GameWebService = &restful.WebService{}
)

func init() {

	id := GameWebService.QueryParameter("playerID", "Identifier of the player").DataType("long")

	GameWebService.
		Path("/v1/game").
		Doc("Starts game").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML).
		ApiVersion("1.0")

	GameWebService.Route(
		GameWebService.GET("/start").Param(id).
			To(startGame))

}
