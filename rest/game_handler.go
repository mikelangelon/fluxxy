package rest

import (
	"github.com/emicklei/go-restful"
	"net/http"
	"strconv"
	"game"
)

func startGame(req *restful.Request, resp *restful.Response) {
	playerId, err := strconv.ParseInt(req.QueryParameter("playerID"), 10, 64)
	if err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	game := game.StartGame(playerId)
	ServeJSON(resp, http.StatusOK, &game)
}
