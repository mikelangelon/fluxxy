package main

import (
	"bol/contract/db"
	"bol/contract/rest"
	"bol/contract/signal"
	"bol/contract/contract_template"
	"github.com/GeertJohan/go.rice"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	log "github.com/inconshreveable/log15"
	"gopkg.in/mgo.v2"
	"net/http"
	"bol/contract/contract"
	"bol/contract/section"
)

//go:generate go get github.com/ReneKroon/go.rice/rice
//go:generate $GOPATH/bin/rice embed-go

func main() {

	mongoDb := MongoDB().DB("")
	defer mongoDb.Session.Close()

	contract_template.TemplateDB = db.CollectionSupplier(mongoDb.Session, contract_template.TemplateCollectionName)
	contract_template.EnsureTemplateIndice()

	contract.ContractDB = db.CollectionSupplier(mongoDb.Session, contract.ContractCollectionName)

	section.SectionDB = db.CollectionSupplier(mongoDb.Session, section.SectionCollectionName)
	// HTTP
	apidocs := "/internal/apidocs/"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", apidocs)
		w.WriteHeader(http.StatusSeeOther)
	})

	http.Handle(apidocs, http.StripPrefix(apidocs, http.FileServer(rice.MustFindBox("../../../../../swagger-dist").HTTPBox())))
	http.HandleFunc("/internal/selfdiagnose.html", rest.HandleSelfdiagnose)

	v1 := rest.V1Handler(*stubsEnabled)
	http.Handle("/v1/", v1)

	webServices := v1.RegisteredWebServices()
	{
		c := restful.NewContainer()
		c.EnableContentEncoding(true)
		swagger.RegisterSwaggerService(swagger.Config{
			WebServices:     webServices,
			ApiPath:         "/internal/apidocs.json",
			SwaggerPath:     "/internal/apidocs-dummy-value/",
			SwaggerFilePath: "swagger-dist",
		}, c)
		http.Handle("/internal/apidocs.json", c)
		http.Handle("/internal/apidocs.json/", c)
	}

	go signal.SignalHandler()
	log.Info("Serving HTTP on interface ", "interface", *httpInterface)
	log.Crit(http.ListenAndServe(*httpInterface, nil).Error())
}

func MongoDB() *mgo.Session {
	s, err := mgo.Dial(*mongoURL)
	if err != nil {
		log.Crit("MongoDB connect", "address", *mongoURL, "error", err)
	}
	log.Info("Using MongoDB ", "address", *mongoURL)

	if *mongoUser != "" {
		creds := &mgo.Credential{Username: *mongoUser, Password: *mongoPass}
		if err := s.Login(creds); err != nil {
			s.Close()
			log.Crit("MongoDB login", "user", *mongoUser, "address", *mongoURL, "error", err)
		}
	}
	return s
}
