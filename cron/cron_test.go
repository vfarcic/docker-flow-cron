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

	s.verifyServicesAreCreated("my-job", 1)
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
		id, _ := exec.Command("/bin/sh", "-c", `docker service ls -q -f label=com.df.cron=true`).CombinedOutput()
		idString := strings.Trim(string(id), "\n")
		if len(id) > 0 {
			out, _ := exec.Command("/bin/sh", "-c", `docker service inspect `+idString).CombinedOutput()
			s.Contains(string(out), `"Condition": "none",`)
			break
		}
		counter++
		if counter >= 100 {
			s.Fail("Service was not created")
			break
		}
		time.Sleep(100 * time.Millisecond)
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
		if counter >= 100 {
			s.Fail("Service was not created")
			break
		}
		time.Sleep(100 * time.Millisecond)
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

// RemoveJob

func (s CronTestSuite) Test_RemoveJob_RemovesService() {
	data := JobData{
		Name:    "my-job",
		Image:   "alpine",
		Command: `echo "Hello Cron!"`,
		Schedule: "@every 1s",
	}

	c, _ := New("unix:///var/run/docker.sock")
	defer func() {
		c.Stop()
		time.Sleep(1 * time.Second)
		s.removeAllServices()
	}()
	c.AddJob(data)
	s.verifyServiceIsCreated()

	c.RemoveJob("my-job")

	counter := 0
	for {
		count := s.getServiceCount("my-job")
		if count == 0 {
			break
		}
		counter++
		if counter >= 50 {
			s.Fail("Services were not removed")
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (s CronTestSuite) Test_RemoveJob_DoesNotRemoveOtherServices() {
	data := JobData{
		Name:    "my-job",
		Image:   "alpine",
		Command: `echo "Hello Cron!"`,
		Schedule: "@every 1s",
	}

	c, _ := New("unix:///var/run/docker.sock")
	defer func() {
		c.Stop()
		time.Sleep(1 * time.Second)
		s.removeAllServices()
	}()
	c.AddJob(data)
	data.Name = "my-other-job"
	c.AddJob(data)
	s.verifyServiceIsCreated()

	c.RemoveJob("my-job")

	before := s.getServiceCount("my-other-job")
	counter := 0
	for {
		after := s.getServiceCount("my-other-job")
		if after > before {
			break
		}
		counter++
		if counter >= 50 {
			s.Fail(fmt.Sprintf("Found %d services. The number should be bigger then %d.", after, before))
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// RescheduleJobs

func (s CronTestSuite) Test_RescheduleJobs_AddsAllJobs() {
	data := JobData{
		Name:    "my-job",
		Image:   "alpine",
		Command: `echo "Hello Cron!"`,
		Schedule: "@every 1s",
	}

	c, _ := New("unix:///var/run/docker.sock")
	defer func() {
		c.Stop()
		time.Sleep(2 * time.Second)
		s.removeAllServices()
	}()
	c.AddJob(data)
	for {
		if s.getServiceCount("my-job") > 0 {
			break
		}
	}
	s.verifyServicesAreCreated("my-job", 1)
	c.Stop()

	c.RescheduleJobs()

	s.verifyServicesAreCreated("my-job", 3)
}

// Util

func (s CronTestSuite) getServiceCount(jobName string) int {
	command := fmt.Sprintf(
		`docker service ls -q -f label=com.df.cron=true -f "label=com.df.cron.name=%s"`,
		jobName,
	)
	out, _ := exec.Command(
		"/bin/sh",
		"-c",
		command,
	).CombinedOutput()
	servicesString := strings.TrimRight(string(out), "\n")
	if len(servicesString) > 0 {
		return len(strings.Split(servicesString, "\n"))
	} else {
		return 0
	}
}

func (s CronTestSuite) verifyServicesAreCreated(serviceName string, replicas int) {
	counter := 0
	for {
		count := s.getServiceCount(serviceName)
		if count >= replicas {
			break
		}
		counter++
		if counter >= 50 {
			s.Fail("Services were not created")
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}

// TODO: Replace with verifyServicesAreCreated
func (s CronTestSuite) verifyServiceIsCreated() {
	counter := 0
	for {
		out, _ := exec.Command(
			"/bin/sh",
			"-c",
			`docker service ls -q -f label=com.df.cron=true -f "label=com.df.cron.name=my-job" -f "label=com.df.cron.schedule=@every 1s"`,
		).CombinedOutput()
		servicesString := strings.TrimRight(string(out), "\n")
		if len(servicesString) > 0 {
			break
		}
		counter++
		if counter >= 15 {
			s.Fail("Service was not created")
			break
		}
		time.Sleep(1 * time.Second)
	}
}

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
