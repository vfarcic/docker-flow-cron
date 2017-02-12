package cron

import (
	"fmt"
	rcron "github.com/robfig/cron"
	"os/exec"
	"strings"
	"../docker"
	"github.com/docker/docker/api/types/swarm"
)

const dockerApiVersion = "v1.24"

type Croner interface {
	AddJob(data JobData) error
	Stop()
	GetJobs() (map[string]JobData, error)
}

type Cron struct {
	Cron *rcron.Cron
	Service docker.Servicer
}

var rCronAddFunc = func(c *rcron.Cron, spec string, cmd func()) error {
	return c.AddFunc(spec, cmd)
}

type JobData struct {
	Name     string
	Image    string
	Command  string
	Schedule string
	Args     []string
}

var New = func(dockerHost string) (Croner, error) {
	service, err := docker.New(dockerHost)
	if err != nil {
		return &Cron{}, err
	}
	c := rcron.New()
	c.Start()
	return &Cron{Cron: c, Service: service}, nil
}

func (c *Cron) AddJob(data JobData) error {
	if data.Args == nil {
		data.Args = []string{}
	}
	if len(data.Name) == 0 {
		return fmt.Errorf("--name is mandatory")
	}
	if len(data.Image) == 0 {
		return fmt.Errorf("image is mandatory")
	}
	cmdPrefix := "docker service create"
	hasRestartCondition := false
	cmdSuffix := ""
	for _, v := range data.Args {
		if strings.HasPrefix(v, "--restart-condition") {
			if strings.Contains(v, "any") {
				return fmt.Errorf("--restart-condition cannot be set to any")
			}
			hasRestartCondition = true
		} else if strings.HasPrefix(v, "--name") {
			return fmt.Errorf("--name argument is not allowed")
		}
		cmdSuffix = fmt.Sprintf("%s %s", cmdSuffix, v)
	}
	if !hasRestartCondition {
		cmdSuffix = fmt.Sprintf("%s --restart-condition none", cmdSuffix)
	}
	cmdSuffix = fmt.Sprintf("%s %s %s", cmdSuffix, data.Image, data.Command)
	cmdLabel := fmt.Sprintf(
		`-l 'com.df.cron.command=%s%s'`,
		cmdPrefix,
		cmdSuffix,
	)
	cmd := fmt.Sprintf(
		`%s -l "com.df.cron=true" -l "com.df.cron.name=%s" -l "com.df.cron.schedule=%s" %s %s`,
		cmdPrefix,
		data.Name,
		data.Schedule,
		cmdLabel,
		strings.Trim(cmdSuffix, " "),
	)
	cronCmd := func() {
		_, err := exec.Command("/bin/sh", "-c", cmd).CombinedOutput()
		if err != nil { // TODO: Test
			fmt.Println(err.Error())
		}
	}
	return rCronAddFunc(c.Cron, data.Schedule, cronCmd)
}

func (c *Cron) GetJobs() (map[string]JobData, error) {
	jobs := map[string]JobData{}
	services, err := c.Service.GetServices("")
	if err != nil {
		return jobs, err
	}
	for _, service := range services {
		command := ""
		for _, v := range service.Spec.TaskTemplate.ContainerSpec.Args {
			if strings.Contains(v, " ") {
				command = fmt.Sprintf(`%s "%s"`, command, v)
			} else {
				command = fmt.Sprintf(`%s %s`, command, v)
			}
		}
		name := service.Spec.Annotations.Labels["com.df.cron.name"]
		jobs[name] = c.getJob(service)
	}
	return jobs, nil
}

func (c *Cron) Stop() {
	c.Cron.Stop()
}

func (c *Cron) getJob(service swarm.Service) JobData {
	command := ""
	for _, v := range service.Spec.TaskTemplate.ContainerSpec.Args {
		if strings.Contains(v, " ") {
			command = fmt.Sprintf(`%s "%s"`, command, v)
		} else {
			command = fmt.Sprintf(`%s %s`, command, v)
		}
	}
	name := service.Spec.Annotations.Labels["com.df.cron.name"]
	return JobData{
		Name:     name,
		Image:    service.Spec.TaskTemplate.ContainerSpec.Image,
		Command:  service.Spec.Annotations.Labels["com.df.cron.command"],
		Schedule: service.Spec.Annotations.Labels["com.df.cron.schedule"],
	}
}
