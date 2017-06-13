package cron

import (
	"../docker"
	"fmt"
	"github.com/docker/docker/api/types/swarm"
	rcron "gopkg.in/robfig/cron.v2"
	"os/exec"
	"strings"
)

const dockerApiVersion = "v1.24"

type Croner interface {
	AddJob(data JobData) error
	Stop()
	GetJobs() (map[string]JobData, error)
	RemoveJob(jobName string) error
	RescheduleJobs() error
}

type Cron struct {
	Cron    *rcron.Cron
	Service docker.Servicer
	Jobs    map[string]rcron.EntryID
}

var rCronAddFunc = func(c *rcron.Cron, spec string, cmd func()) (rcron.EntryID, error) {
	return c.AddFunc(spec, cmd)
}

type JobData struct {
	Name           string   `json:"name"`
	ServiceName    string   `json:"servicename"`
	Image          string   `json:"image"`
	Command        string   `json:"command"`
	Schedule       string   `json:"schedule"`
	Args           []string `json:"args"`
	Created        bool     `json:"created`
}

var New = func(dockerHost string) (Croner, error) {
	service, err := docker.New(dockerHost)
	if err != nil {
		return &Cron{}, err
	}
	c := rcron.New()
	c.Start()
	return &Cron{Cron: c, Service: service, Jobs: map[string]rcron.EntryID{}}, nil
}

func (c *Cron) AddJob(data JobData) error {
	fmt.Println("Scheduling", data.Name)
	if data.Args == nil {
		data.Args = []string{}
	}
	if len(data.Name) == 0 {
		return fmt.Errorf("name is mandatory")
	}
	if len(data.Image) == 0 {
		return fmt.Errorf("image is mandatory")
	}
	/*if len(data.Schedule) == 0 {
		return fmt.Errorf("schedule is mandatory")
	}*/

	serviceName := data.Name
	if data.ServiceName != "" {
		serviceName = data.ServiceName
	}

	if !data.Created{
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
			`%s -l "com.df.cron=true" -l "com.df.cron.name=%s" -l "com.df.cron.schedule=%s" --name %s --replicas 0 %s %s`,
			cmdPrefix,
			data.Name,
			data.Schedule,
			serviceName,
			cmdLabel,
			strings.Trim(cmdSuffix, " "),
		)

		fmt.Println("Executing command:", cmd)

		_, err := exec.Command("/bin/sh", "-c", cmd).CombinedOutput()
		if err != nil { // TODO: Test
		       fmt.Println("Could not execute command: ", cmd, err.Error())
		}
	}

	cronCmd := func() {
		scale := fmt.Sprintf(`docker service scale %s=1`, serviceName)
		fmt.Println(scale)
		_, err := exec.Command("/bin/sh", "-c", scale).CombinedOutput()
		if err != nil { // TODO: Test
			fmt.Println("Could not execute command: ", scale)
		}
	}
	entryId, err := rCronAddFunc(c.Cron, data.Schedule, cronCmd)
	c.Jobs[data.Name] = entryId
	return err
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

func (c *Cron) RemoveJob(jobName string) error {
	fmt.Println("Removing job", jobName)
	c.Cron.Remove(c.Jobs[jobName])
	if err := c.Service.RemoveServices(jobName); err != nil {
		return err
	}
	return nil
}

func (c *Cron) RescheduleJobs() error {
	fmt.Println("Rescheduling jobs")
	jobs, err := c.GetJobs()
	if err != nil {
		return err
	}
	for _, job := range jobs {
		c.AddJob(job)
	}
	c.Cron.Start()
	return nil
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
		ServiceName: service.Spec.Name,
		Image:    service.Spec.TaskTemplate.ContainerSpec.Image,
		Command:  service.Spec.Annotations.Labels["com.df.cron.command"],
		Schedule: service.Spec.Annotations.Labels["com.df.cron.schedule"],
	}
}
