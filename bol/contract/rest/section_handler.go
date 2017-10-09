package rest

import (
	"github.com/emicklei/go-restful"
	"net/http"
	"bol/contract/section"
	"strconv"
	"github.com/pascaldekloe/goe/rest"
)

func createSection(req *restful.Request, resp *restful.Response) {
	entity := section.Section{}
	err := req.ReadEntity(&entity)
	if err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}

	switch err := section.Insert(entity); {
	case err != nil:
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
}

func updateSection(req *restful.Request, resp *restful.Response) {
	sectionId, err := strconv.ParseInt(req.PathParameter("id"), 10, 64)
	if err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}

	entity := section.Section{}
	entity.ID = sectionId
	err = req.ReadEntity(&entity)
	if err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}

	switch err := section.Update(entity); {
	case err != nil:
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
}

func retrieveSection(req *restful.Request, resp *restful.Response) {
	sectionId, err := strconv.ParseInt(req.PathParameter("id"), 10, 64)
	if err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}
	section, err := section.Lookup(sectionId)
	if err != nil {
		resp.WriteError(http.StatusInternalServerError, err)
		return
	}
	rest.ServeJSON(resp, http.StatusOK, section)
}
