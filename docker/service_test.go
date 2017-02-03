package docker

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"os/exec"
	"testing"
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
	s.createTestService("util-1", "-l com.df.cron=true")
	s.createTestService("util-2", "-l com.df.test=true")
	defer s.removeAllServices()
	services, _ := New("unix:///var/run/docker.sock")

	actual, _ := services.GetServices()

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

	_, err := services.GetServices()

	s.Error(err)
}

// Util

func (s *ServiceTestSuite) createTestService(name, args string) {
	cmd := fmt.Sprintf(
		`docker service create --name %s %s alpine sleep 10000000`,
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
