package contract

import (
	"bytes"
	"fmt"
	diag "github.com/emicklei/go-selfdiagnose"
	"github.com/emicklei/go-selfdiagnose/task"
	"os"
	"runtime"
	"strings"
	"time"
)

var Timestamp, Commithash, Version, Goversion string

func NewDiagRegistry() *diag.Registry {
	reg := new(diag.Registry)

	id := new(IDCheck)
	id.SetComment("ID")
	reg.Register(id)

	proc := new(ProcCheck)
	proc.SetComment("Processing")
	reg.Register(proc)

	mem := new(MemCheck)
	mem.SetComment("Memory")
	reg.Register(mem)

	health := new(HealthCheck)
	health.SetComment("Consul Health")
	health.SetTimeout(2 * time.Second)
	reg.Register(health)

	service := new(ServiceCheck)
	service.SetComment("Consul Service")
	service.SetTimeout(2 * time.Second)
	reg.Register(service)

	kevlar := new(KevlarCheck)
	kevlar.SetComment("Configuration")
	reg.Register(kevlar)

	timezone := new(TimeZoneCheck)
	timezone.SetComment("Default time on machine")
	reg.Register(timezone)

	env := new(EnvCheck)
	env.SetComment("Environment")
	reg.Register(env)

	var build task.ReportBuildAndDate
	build.Version = Version
	build.Date = Timestamp
	reg.Register(build)

	commit := new(CommitCheck)
	commit.SetComment("Git Commit")
	reg.Register(commit)

	goVersion := new(GoVersionCheck)
	goVersion.SetComment("Go Version")
	reg.Register(goVersion)

	cache := new(CacheCheck)
	cache.SetComment("Cache Check")
	reg.Register(cache)

	return reg
}

func SetupSelfdiagnose(build task.ReportBuildAndDate) {
	diag.AddInternalHandlers()
	diag.Register(build)
	diag.Register(task.ReportHttpRequest{})
	diag.Register(task.ReportHostname{})
	diag.Register(task.ReportCPU())
}

type TimeZoneCheck struct{ diag.BasicTask }

func (c *TimeZoneCheck) Run(ctx *diag.Context, result *diag.Result) {
	result.Reason = fmt.Sprintf("%s", time.Now().Format(time.UnixDate))
	result.Passed = true
}

type CacheCheck struct{ diag.BasicTask }

func (c *CacheCheck) Run(ctx *diag.Context, result *diag.Result) {
	result.Passed = true
}

type CommitCheck struct{ diag.BasicTask }

func (c *CommitCheck) Run(ctx *diag.Context, result *diag.Result) {
	result.Reason = fmt.Sprintf("%s", Commithash)
	result.Passed = true
}

type GoVersionCheck struct{ diag.BasicTask }

func (c *GoVersionCheck) Run(ctx *diag.Context, result *diag.Result) {
	result.Reason = fmt.Sprintf("%s", Goversion)
	result.Passed = true
}

type IDCheck struct{ diag.BasicTask }

func (c *IDCheck) Run(ctx *diag.Context, result *diag.Result) {
	result.Reason = fmt.Sprintf("%s uptime %s", ConfLabel, time.Since(Started))
	result.Passed = true
}

type ProcCheck struct{ diag.BasicTask }

func (c *ProcCheck) Run(ctx *diag.Context, result *diag.Result) {
	result.Reason = fmt.Sprintf("%d Go routines for %d CPUs - %d golang threads", runtime.NumGoroutine(), runtime.NumCPU(), runtime.GOMAXPROCS(-1))
	result.Passed = runtime.NumGoroutine() < 1000
	if runtime.NumGoroutine() < 2000 {
		result.Severity = diag.SeverityWarning
	} else if runtime.NumGoroutine() >= 2000 {
		result.Severity = diag.SeverityCritical
	}
}

type MemCheck struct{ diag.BasicTask }

func (c *MemCheck) Run(ctx *diag.Context, result *diag.Result) {
	mem := new(runtime.MemStats)
	runtime.ReadMemStats(mem)

	const mib = 1024 * 1024

	result.Passed = mem.Sys < 1536*mib
	result.Reason = fmt.Sprintf("system %dMiB, allocated %dMiB", mem.Sys/mib, mem.Alloc/mib)

	if mem.Sys < 2048*mib {
		result.Severity = diag.SeverityWarning

	} else if mem.Sys >= 2048*mib {
		result.Severity = diag.SeverityCritical
	}
}

type HealthCheck struct{ diag.BasicTask }

func (c *HealthCheck) Run(ctx *diag.Context, result *diag.Result) {
	if ConsulClient == nil {
		result.Reason = "Consul not configured"
		return
	}

	checks, err := ConsulClient.Agent().Checks()
	if err != nil {
		result.Reason = err.Error()
		return
	}

	buf := new(bytes.Buffer)
	result.Passed = true
	for _, c := range checks {
		fmt.Fprintf(buf, "%s: %q for service %q %s<br>\n", c.Node, c.Name, c.ServiceName, c.Status)
		if c.Status != "passing" {
			result.Passed = false
		}
	}
	result.Reason = buf.String()
}

type ServiceCheck struct{ diag.BasicTask }

func (c *ServiceCheck) Run(ctx *diag.Context, result *diag.Result) {
	if ConsulClient == nil {
		result.Reason = "Consul not configured"
		return
	}

	services, err := ConsulClient.Agent().Services()
	if err != nil {
		result.Reason = err.Error()
		return
	}

	buf := new(bytes.Buffer)
	for _, s := range services {
		fmt.Fprintf(buf, "//%s:%d/: %q %q<br>\n", s.Address, s.Port, s.Service, s.Tags)
	}
	result.Passed = true
	result.Reason = buf.String()
}

type KevlarCheck struct{ diag.BasicTask }

func (c *KevlarCheck) Run(ctx *diag.Context, result *diag.Result) {
	buf := new(bytes.Buffer)

	for _, key := range Conf.Keys() {
		buf.WriteString(key)
		buf.WriteByte('=')
		if strings.HasSuffix(key, "*") {
			buf.WriteString("_hidden_")
		} else {
			v, _ := Conf.Get(key)
			buf.WriteString(v)
		}
		buf.WriteString("<br>\n")
	}

	result.Reason = buf.String()
	result.Passed = true
}

type EnvCheck struct{ diag.BasicTask }

func (c *EnvCheck) Run(ctx *diag.Context, result *diag.Result) {
	result.Reason = strings.Join(os.Environ(), "<br>\n")
	result.Passed = true
}
