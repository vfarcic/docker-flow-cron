package cron

import (
	"fmt"
	rcron "github.com/robfig/cron"
	"os/exec"
)

type Croner interface {
	AddJob(data JobData) error
	Stop()
}

type Cron struct {
	Cron *rcron.Cron
}

var rCronAddFunc = func(c *rcron.Cron, spec string, cmd func()) error {
	return c.AddFunc(spec, cmd)
}

type JobData struct {
	Name    string
	Image   string
	Command string
	Schedule string
	Params  map[string]string
}

var New = func() Croner {
	c := rcron.New()
	c.Start()
	return &Cron{
		Cron: c,
	}
}

func (c *Cron) AddJob(data JobData) error {
	if data.Params == nil {
		data.Params = map[string]string{}
	}
	if len(data.Name) == 0 {
		return fmt.Errorf("--name is mandatory")
	}
	if len(data.Image) == 0 {
		return fmt.Errorf("image is mandatory")
	}
	cmd := fmt.Sprintf(
		`docker service create -l "com.df.cron=true" -l "com.df.cron.name=%s" -l "com.df.cron.schedule=%s"`,
		data.Name,
		data.Schedule,
	)
	if _, ok := data.Params["--restart-condition"]; !ok {
		data.Params["--restart-condition"] = "none"
	}
	for k, v := range data.Params {
		if k == "--restart-condition" && v == "any" {
			return fmt.Errorf("--restart-condition cannot be set to any")
		}
		cmd = fmt.Sprintf("%s %s %s", cmd, k, v)
	}
	cmd = fmt.Sprintf("%s %s %s", cmd, data.Image, data.Command)
	cronCmd := func() {
		_, err := exec.Command("/bin/sh", "-c", cmd).CombinedOutput()
		if err != nil { // TODO: Test
			fmt.Println(err.Error())
		}
	}
	return rCronAddFunc(c.Cron, data.Schedule, cronCmd)
}

func (c *Cron) Stop() {
	c.Cron.Stop()
}
