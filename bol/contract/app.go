package contract

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	log "github.com/inconshreveable/log15"

	consul "github.com/hashicorp/consul/api"
	"github.com/magiconair/properties"
	kevlar "stash.tools.bol.com/kev/kevlar-go-client.git"
)

var (
	// Env holds the applicable environment label.
	ConfEnv = os.Getenv("ENV")
	// ConfLabel is the identification for Kevlar.
	ConfLabel = "contract"
	// ConfVersion is the key set version for Kevlar.
	ConfVersion = "1"

	// Conf holds the configuration properties.
	Conf *properties.Properties
)

var ConsulClient *consul.Client

var HTTPInterface string

// Started records the uptime.
var Started = time.Now()

func LoadConf() *properties.Properties {
	if ConfEnv == "" {
		log.Warn("Default to dev due missing ENV environment variable")
		ConfEnv = "dev"
	}

	baseURL := "http://shd-kevlar-app-vip.bolcom.net"
	if ConfEnv == "dev" {
		log.Warn("Running in development mode")
		baseURL = "https://kevlar-dev.tools.bol.com"
	}

	if ConfEnv == "jenkins" {
		log.Warn("Running on jenkins, skipping kevlar")
		return properties.NewProperties()
	}

	if ConfEnv == "build" {
		log.Warn("Running a build, skipping kevlar")
		return properties.NewProperties()
	}
	return kevlar.MustResolve(baseURL, ConfLabel, ConfEnv, ConfVersion)
}

func Registration() (*consul.AgentServiceRegistration, []*consul.AgentCheckRegistration) {
	addr, portStr, err := net.SplitHostPort(HTTPInterface)
	if err != nil {
		log.Error("Consul resolve port", "interface", HTTPInterface, "error", err)
	}
	if addr == "" {
		if addr, err = os.Hostname(); err != nil {
			log.Error("Resolve hostname: ", "error", err)
		}
	}
	port, err := strconv.ParseInt(portStr, 10, 16)
	if err != nil {
		log.Error("Consul malformed port number", "interface", HTTPInterface, "error", err)
	}

	serviceReg := &consul.AgentServiceRegistration{
		ID:      "contract",
		Name:    "contract",
		Tags:    []string{"pdfMaker", "rest"},
		Port:    int(port),
		Address: addr,
	}

	checkRegs := []*consul.AgentCheckRegistration{
		{
			ID:        "self-diagnostics",
			Name:      "self-diagnostics",
			Notes:     "HTTP lookup of the page.",
			ServiceID: "contract",
			AgentServiceCheck: consul.AgentServiceCheck{
				HTTP:     fmt.Sprintf("http://%s:%d/internal/selfdiagnose.html", addr, port),
				Interval: "10s",
				Timeout:  "1s",
			},
		},
	}

	for _, c := range checkRegs {
		serviceReg.Checks = append(serviceReg.Checks, &c.AgentServiceCheck)
	}

	return serviceReg, checkRegs
}
