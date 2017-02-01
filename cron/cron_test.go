package cron

import (
	rcron "github.com/robfig/cron"
	"github.com/stretchr/testify/suite"
	"os/exec"
	"strings"
	"testing"
	"time"
)

type CronTestSuite struct {
	suite.Suite
}

func (s *CronTestSuite) SetupTest() {
}

func TestCronUnitTestSuite(t *testing.T) {
	s := new(CronTestSuite)
	suite.Run(t, s)
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
		Name:    "my-job",
		Image:   "alpine",
		Schedule: "@yearly",
		Params:  map[string]string{"--restart-condition": "any"},
		Command: `echo "Hello Cron!"`,
	}
	c := New()

	err := c.AddJob(data)

	s.Error(err)
}

func (s CronTestSuite) Test_AddJob_AddsRestartConditionNone_WhenNotSet() {
	data := JobData{
		Name:    "my-job",
		Image:   "alpine",
		Params:  map[string]string{},
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

func (s CronTestSuite) Test_AddJob_ThrowsAnError_WhenImageIsEmpty() {
	data := JobData{
		Name:    "my-job",
		Image:   "",
		Schedule: "@yearly",
		Command: `echo "Hello Cron!"`,
	}
	c := New()

	err := c.AddJob(data)

	s.Error(err)
}

func (s CronTestSuite) Test_AddJob_ThrowsAnError_WhenNameIsEmpty() {
	data := JobData{
		Image:   "my-image",
		Schedule: "@yearly",
		Command: `echo "Hello Cron!"`,
	}
	c := New()

	err := c.AddJob(data)

	s.Error(err)
}

func (s CronTestSuite) addJob1s(d JobData) {
	c := New()
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
