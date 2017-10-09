package rest

import (
	"github.com/emicklei/go-restful"
	"net/http"
	"bol/contract/contract"
	"github.com/pascaldekloe/goe/rest"
	"strconv"
)

func createContract(req *restful.Request, resp *restful.Response) {
	entity := contract.Contract{}
	err := req.ReadEntity(&entity)
	if err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}

	switch err := contract.Insert(entity); {
	case err != nil:
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
}

func updateContract(req *restful.Request, resp *restful.Response) {
	contractId, err := strconv.ParseInt(req.PathParameter("id"), 10, 64)
	if err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}

	entity := contract.Contract{}
	entity.ID = contractId
	err = req.ReadEntity(&entity)
	if err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}

	switch err := contract.Update(entity); {
	case err != nil:
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
}

func retrieveContract(req *restful.Request, resp *restful.Response) {
	contractId, err := strconv.ParseInt(req.PathParameter("id"), 10, 64)
	if err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	template, err := contract.Lookup(contractId)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	rest.ServeJSON(resp, http.StatusOK, template)
}
