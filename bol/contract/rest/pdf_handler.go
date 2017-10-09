package rest

import (
	"bol/contract/pdf"
	"bytes"
	"github.com/emicklei/go-restful"
	"net/http"
	"strconv"
)

func generatePDF(req *restful.Request, resp *restful.Response) {
	document := pdf.GeneratePdf()
	buf := make([]byte, 0)
	w := bytes.NewBuffer(buf)
	document.Output(w)

	writer := resp.ResponseWriter
	writer.Header().Set("Content-Type", "application/pdf")
	writer.Header().Add("Content-Disposition", "attachment; filename=test.pdf")
	writer.Header().Set("Content-Length", strconv.Itoa(len(w.Bytes())))
	if _, err := writer.Write(w.Bytes()); err != nil {
		http.Error(resp, "Cannot encode pdf"+err.Error(), http.StatusInternalServerError)
		return
	}
}

func generatePDFFromTemplate(req *restful.Request, resp *restful.Response) {
	templateId, err := strconv.ParseInt(req.PathParameter("id"), 10, 64)
	if err != nil {
		resp.WriteError(http.StatusBadRequest, err)
		return
	}

	document := pdf.GeneratePdfFromContract(templateId)
	buf := make([]byte, 0)
	w := bytes.NewBuffer(buf)
	document.Output(w)

	writer := resp.ResponseWriter
	writer.Header().Set("Content-Type", "application/pdf")
	writer.Header().Add("Content-Disposition", "attachment; filename=test.pdf")
	writer.Header().Set("Content-Length", strconv.Itoa(len(w.Bytes())))
	if _, err := writer.Write(w.Bytes()); err != nil {
		http.Error(resp, "Cannot encode pdf"+err.Error(), http.StatusInternalServerError)
		return
	}
}
