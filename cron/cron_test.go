package cron

import (
	rcron "github.com/robfig/cron"
	"github.com/stretchr/testify/suite"
	"os/exec"
	"strings"
	"testing"
	"time"
	"fmt"
	"github.com/docker/docker/api/types/swarm"
)

type CronTestSuite struct {
	suite.Suite
	Service        ServicerMock
}

func (s *CronTestSuite) SetupTest() {
	s.Service = ServicerMock{
		GetServicesMock: func(jobName string) ([]swarm.Service, error) {
			return []swarm.Service{}, nil
		},
		GetTasksMock: func(jobName string) ([]swarm.Task, error) {
			return []swarm.Task{}, nil
		},
	}
}

func TestCronUnitTestSuite(t *testing.T) {
	s := new(CronTestSuite)
	suite.Run(t, s)
}

// New

func (s *CronTestSuite) Test_New_ReturnsError_WhenDockerClientFails() {
	_, err := New("this-is-not-a-socket")

	s.Error(err)
}

// AddJob

func (s CronTestSuite) Test_AddJob_InvokesRCronAddFuncWithSpec() {
	rCronAddFuncOrig := rCronAddFunc
	defer func() { rCronAddFunc = rCronAddFuncOrig }()
	actualSpec := ""
	rCronAddFunc = func(c *rcron.Cron, spec string, cmd func()) error {
		actualSpec = spec
		return nil
	}
	data := JobData{Image: "my-image", Name: "my-job", Schedule: "@yearly"}
	c := Cron{
		Cron: rcron.New(),
	}

	c.AddJob(data)

	s.Equal(data.Schedule, actualSpec)
}

func (s CronTestSuite) Test_AddJob_CreatesService() {
	data := JobData{
		Name:    "my-job",
		Image:   "alpine",
		Command: `echo "Hello Cron!"`,
	}
	defer s.removeAllServices()

	s.addJob1s(data)

	counter := 0
	for {
		time.Sleep(1 * time.Second)
		out, _ := exec.Command(
			"/bin/sh",
			"-c",
			`docker service ls -q -f label=com.df.cron=true -f "label=com.df.cron.name=my-job" -f "label=com.df.cron.schedule=@every 1s"`,
		).CombinedOutput()
		if len(out) > 0 {
			break
		}
		counter++
		if counter >= 15 {
			s.Fail("Service was not created")
			break
		}
	}
}

func (s CronTestSuite) Test_AddJob_ThrowsAnError_WhenRestartConditionIsSetToAny() {
	data := JobData{
		Name:     "my-job",
		Image:    "alpine",
		Schedule: "@yearly",
		Args:     []string{"--restart-condition any"},
		Command:  `echo "Hello Cron!"`,
	}
	c, _ := New("unix:///var/run/docker.sock")

	err := c.AddJob(data)

	s.Error(err)
}

func (s CronTestSuite) Test_AddJob_AddsRestartConditionNone_WhenNotSet() {
	data := JobData{
		Name:    "my-job",
		Image:   "alpine",
		Args:    []string{},
		Command: `echo "Hello Cron!"`,
	}
	s.addJob1s(data)
	defer s.removeAllServices()

	counter := 0
	for {
		time.Sleep(1 * time.Second)
		id, _ := exec.Command("/bin/sh", "-c", `docker service ls -q -f label=com.df.cron=true`).CombinedOutput()
		idString := strings.Trim(string(id), "\n")
		if len(id) > 0 {
			out, _ := exec.Command("/bin/sh", "-c", `docker service inspect `+idString).CombinedOutput()
			s.Contains(string(out), `"Condition": "none",`)
			break
		}
		counter++
		if counter >= 15 {
			s.Fail("Service was not created")
			break
		}
	}

}

func (s CronTestSuite) Test_AddJob_ThrowsAnError_WhenNameArgumentIsSet() {
	data := JobData{
		Name:     "my-job",
		Image:    "alpine",
		Schedule: "@yearly",
		Args:     []string{"--name some-name"},
		Command:  `echo "Hello Cron!"`,
	}
	c, _ := New("unix:///var/run/docker.sock")

	err := c.AddJob(data)

	s.Error(err)
}

func (s CronTestSuite) Test_AddJob_AddsCommandLabel() {
	data := JobData{
		Name:    "my-job",
		Image:   "alpine",
		Args:    []string{},
		Command: `echo "Hello Cron!"`,
	}
	s.addJob1s(data)
	defer s.removeAllServices()

	counter := 0
	for {
		time.Sleep(1 * time.Second)
		id, _ := exec.Command("/bin/sh", "-c", `docker service ls -q -f label=com.df.cron=true`).CombinedOutput()
		idString := strings.Trim(string(id), "\n")
		if len(id) > 0 {
			out, _ := exec.Command("/bin/sh", "-c", `docker service inspect `+idString).CombinedOutput()
			s.Contains(
				string(out),
				`"com.df.cron.command": "docker service create --restart-condition none alpine echo \"Hello Cron!\""`,
			)
			break
		}
		counter++
		if counter >= 15 {
			s.Fail("Service was not created")
			break
		}
	}

}

func (s CronTestSuite) Test_AddJob_ThrowsAnError_WhenImageIsEmpty() {
	data := JobData{
		Name:     "my-job",
		Image:    "",
		Schedule: "@yearly",
		Command:  `echo "Hello Cron!"`,
	}
	c, _ := New("unix:///var/run/docker.sock")

	err := c.AddJob(data)

	s.Error(err)
}

func (s CronTestSuite) Test_AddJob_ThrowsAnError_WhenNameIsEmpty() {
	data := JobData{
		Image:    "my-image",
		Schedule: "@yearly",
		Command:  `echo "Hello Cron!"`,
	}
	c, _ := New("unix:///var/run/docker.sock")

	err := c.AddJob(data)

	s.Error(err)
}

// GetJobs

func (s CronTestSuite) Test_GetJobs_ReturnsListOfJobs() {
	expected := map[string]JobData{}
	defer s.removeAllServices()
	for i := 1; i <= 3; i++ {
		name := fmt.Sprintf("my-job-%d", i)
		cmd := fmt.Sprintf(
			`docker service create \
    -l 'com.df.cron=true' \
    -l 'com.df.cron.name=%s' \
    -l 'com.df.cron.schedule=@every 1s' \
    -l 'com.df.cron.command=docker service create --restart-condition none alpine echo "Hello World!"' \
    --constraint "node.labels.env != does-not-exist" \
    --container-label 'container=label' \
    --restart-condition none \
    alpine:3.5@sha256:dfbd4a3a8ebca874ebd2474f044a0b33600d4523d03b0df76e5c5986cb02d7e8 \
    echo "Hello world!"`,
			name,
		)
		exec.Command("/bin/sh", "-c", cmd).CombinedOutput()
		expected[name] = JobData{
			Name:     name,
			Image:    "alpine:3.5@sha256:dfbd4a3a8ebca874ebd2474f044a0b33600d4523d03b0df76e5c5986cb02d7e8",
			Command:  `docker service create --restart-condition none alpine echo "Hello World!"`,
			Schedule: "@every 1s",
		}
	}

	c, _ := New("unix:///var/run/docker.sock")
	actual, _ := c.GetJobs()

	s.Equal(expected, actual)
}

func (s *CronTestSuite) Test_GetJobs_ReturnsError_WhenGetServicesFail() {
	message := "This is an error"
	mock := ServicerMock{
		GetServicesMock: func(jobName string) ([]swarm.Service, error) {
			return []swarm.Service{}, fmt.Errorf(message)
		},
	}

	c := Cron{Cron: rcron.New(), Service: mock}
	_, err := c.GetJobs()

	s.Error(err)
}

// Util

func (s CronTestSuite) addJob1s(d JobData) {
	c, _ := New("unix:///var/run/docker.sock")
	d.Schedule = "@every 1s"
	c.AddJob(d)
	time.Sleep(1 * time.Second)
	c.Stop()
}

func (s CronTestSuite) removeAllServices() {
	exec.Command(
		"/bin/sh",
		"-c",
		`docker service rm $(docker service ls -q -f label=com.df.cron=true)`,
	).CombinedOutput()
}

type ServicerMock struct {
	GetServicesMock func(jobName string) ([]swarm.Service, error)
	GetTasksMock    func(jobName string) ([]swarm.Task, error)
	RemoveServicesMock func(jobName string) error
}

func (m ServicerMock) GetServices(jobName string) ([]swarm.Service, error) {
	return m.GetServicesMock(jobName)
}

func (m ServicerMock) GetTasks(jobName string) ([]swarm.Task, error) {
	return m.GetTasksMock(jobName)
}

func (m ServicerMock) RemoveServices(jobName string) error {
	return m.RemoveServicesMock(jobName)
}
