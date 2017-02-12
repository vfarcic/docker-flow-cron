package docker

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"os/exec"
	"testing"
	"time"
)

type ServiceTestSuite struct {
	suite.Suite
}

func (s *ServiceTestSuite) SetupTest() {
}

func TestServiceUnitTestSuite(t *testing.T) {
	s := new(ServiceTestSuite)
	suite.Run(t, s)
}

// New

func (s *ServiceTestSuite) Test_New_ReturnsError_WhenClientFails() {
	expected := "this-is-a-host"

	_, err := New(expected)

	s.Error(err)
}

func (s *ServiceTestSuite) Test_New_SetsClient() {
	expected := "unix:///var/run/docker.sock"

	srv, _ := New(expected)

	s.NotNil(srv.Client)
}

// GetServices

func (s *ServiceTestSuite) Test_GetServices_ReturnsServices() {
	defer s.removeAllServices()
	s.createTestService("util-1", "-l com.df.cron=true")
	s.createTestService("util-2", "-l com.df.test=true")
	services, _ := New("unix:///var/run/docker.sock")

	actual, _ := services.GetServices("")

	if len(actual) == 0 {
		s.Fail("No services found")
	} else {
		s.Equal(1, len(actual))
		s.Equal("util-1", actual[0].Spec.Name)
		s.Equal("true", actual[0].Spec.Labels["com.df.cron"])
	}
}

func (s *ServiceTestSuite) Test_GetServices_ReturnsError_WhenServiceListFails() {
	services, _ := New("unix:///this/socket/does/not/exist")

	_, err := services.GetServices("")

	s.Error(err)
}

func (s *ServiceTestSuite) Test_GetServices_ReturnsFilteredServices() {
	defer s.removeAllServices()
	s.createTestService("util-3", "-l com.df.cron.name=my-job -l com.df.cron=true")
	s.createTestService("util-4", "-l com.df.cron.name=my-job -l com.df.cron=true")
	s.createTestService("util-5", "-l com.df.cron.name=some-other-job -l com.df.cron=true")
	services, _ := New("unix:///var/run/docker.sock")

	actual, _ := services.GetServices("my-job")

	if len(actual) == 0 {
		s.Fail("No services found")
	} else {
		s.Equal(2, len(actual))
	}
}

// GetTasks

func (s *ServiceTestSuite) Test_GetTasks_ReturnsFilteredTasks() {
	defer s.removeAllServices()
	s.createTestService("util-6", "-l com.df.cron.name=my-job -l com.df.cron=true")
	s.createTestService("util-7", "-l com.df.cron.name=my-job -l com.df.cron=true")
	s.createTestService("util-8", "-l com.df.cron.name=some-other-job -l com.df.cron=true")
	services, _ := New("unix:///var/run/docker.sock")
	actual := 0

	for i := 0; i < 100; i++ {
		t, _ := services.GetTasks("my-job")
		actual = len(t)
		if actual == 2 {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}

	s.Equal(2, actual)
}

func (s *ServiceTestSuite) Test_GetTasks_ReturnsError_WhenServiceListFails() {
	services, _ := New("unix:///this/socket/does/not/exist")

	_, err := services.GetTasks("some-job")

	s.Error(err)
}

// RemoveServices

func (s *ServiceTestSuite) Test_RemoveServices_RemovesAllServicesRelatedToACronJob() {
	defer s.removeAllServices()
	s.createTestService("util-9", "-l com.df.cron.name=my-job -l com.df.cron=true")
	s.createTestService("util-10", "-l com.df.cron.name=my-job -l com.df.cron=true")
	s.createTestService("util-11", "-l com.df.cron.name=some-other-job -l com.df.cron=true")
	services, _ := New("unix:///var/run/docker.sock")

	services.RemoveServices("my-job")

	myJobServices, _ := services.GetServices("my-job")
	s.Equal(0, len(myJobServices))

	someOtherJobServices, _ := services.GetServices("some-other-job")
	s.Equal(1, len(someOtherJobServices))
}

func (s *ServiceTestSuite) Test_RemoveServices_ReturnsAnError_WhenClientFails() {
	defer s.removeAllServices()
	services, _ := New("unix:///this/socket/does/not/exist")

	err := services.RemoveServices("my-job")

	s.Error(err)
}

// Util

func (s *ServiceTestSuite) createTestService(name, args string) {
	cmd := fmt.Sprintf(
		`docker service create --name %s --restart-condition none %s alpine sleep 10000000`,
		name,
		args,
	)
	exec.Command("/bin/sh", "-c", cmd).CombinedOutput()
}

func (s *ServiceTestSuite) removeAllServices() {
	exec.Command(
		"/bin/sh",
		"-c",
		`docker service rm $(docker service ls -q -f label=com.df.cron=true)`,
	).CombinedOutput()
	exec.Command(
		"/bin/sh",
		"-c",
		`docker service rm $(docker service ls -q -f label=com.df.test=true)`,
	).CombinedOutput()
}
