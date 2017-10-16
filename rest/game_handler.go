package rest

import (
	"github.com/emicklei/go-restful"
	"net/http"
	"strconv"
	"mainmodel"
)

func startGame(req *restful.Request, resp *restful.Response){
	playerId, err := strconv.ParseInt(req.QueryParameter("playerID"), 10, 64)
	if err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	ServeJSON(resp, http.StatusOK, &mainmodel.Game{ID: 1,
		Players: []*mainmodel.Player{{ID:playerId},{ID:playerId+1},{ID:playerId+2},{ID:playerId+3}}})
}
