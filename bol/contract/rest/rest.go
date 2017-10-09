// Package rest provides the HTTP API.
package rest

import (
	"bol/contract"
	"github.com/emicklei/go-restful"
	diag "github.com/emicklei/go-selfdiagnose"
	"github.com/pascaldekloe/goe/metrics"
	"net/http"
)

var Metrics metrics.Register = metrics.NewDummy()

func init() {

}

func V1Handler(doStub bool) *restful.Container {
	c := restful.NewContainer()
	c.Add(PdfWebService)
	c.Add(ContractWebService)
	c.Add(TemplateWebService)
	c.Add(SectionWebService)
	return c
}

func HandleSelfdiagnose(w http.ResponseWriter, r *http.Request) {
	var reporter diag.Reporter
	switch r.URL.Query().Get("format") {
	case "", "html":
		reporter = diag.HtmlReporter{Writer: w}
	case "json":
		reporter = diag.JSONReporter{Writer: w}
	case "xml":
		reporter = diag.XMLReporter{Writer: w}
	default:
		http.Error(w, "unsupported format", http.StatusBadRequest)
		return
	}

	contract.NewDiagRegistry().Run(reporter)
}
