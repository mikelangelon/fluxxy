package rest

import (
	"bol/contract/contract_template"
	"github.com/emicklei/go-restful"
	"github.com/pascaldekloe/goe/rest"
	"net/http"
	"strconv"
)

func createTemplate(req *restful.Request, resp *restful.Response) {
	entity := contract_template.Template{}
	err := req.ReadEntity(&entity)
	if err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}

	switch err := contract_template.Insert(entity); {
	case err != nil:
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
}

func updateTemplate(req *restful.Request, resp *restful.Response) {
	templateId, err := strconv.ParseInt(req.PathParameter("id"), 10, 64)
	if err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}

	entity := contract_template.Template{}
	entity.ID = templateId
	err = req.ReadEntity(&entity)
	if err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}

	switch err := contract_template.Update(entity); {
	case err != nil:
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
}

func retrieveTemplate(req *restful.Request, resp *restful.Response) {
	templateId, err := strconv.ParseInt(req.PathParameter("id"), 10, 64)
	if err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	template, err := contract_template.Lookup(templateId)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	rest.ServeJSON(resp, http.StatusOK, template)
}
