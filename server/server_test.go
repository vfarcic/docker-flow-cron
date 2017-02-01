package server

import (
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
	"fmt"
	"net/http"
	"github.com/stretchr/testify/mock"
	"../cron"
	"strings"
	"encoding/json"
)

type ServerTestSuite struct {
	suite.Suite
	ResponseWriter     *ResponseWriterMock
}

func (s *ServerTestSuite) SetupTest() {
	s.ResponseWriter = getResponseWriterMock()
}

func TestCronUnitTestSuite(t *testing.T) {
	s := new(ServerTestSuite)
	suite.Run(t, s)
}

// Execute

func (s *ServerTestSuite) Test_Execute_InvokesHTTPListenAndServe() {
	serve := New("myIp", "1234")
	var actual string
	expected := fmt.Sprintf("%s:%s", serve.IP, serve.Port)
	httpListenAndServe = func(addr string, handler http.Handler) error {
		actual = addr
		return nil
	}

	serve.Execute()
	time.Sleep(1 * time.Millisecond)

	s.Equal(expected, actual)
}

func (s *ServerTestSuite) Test_Execute_ReturnsError_WhenHTTPListenAndServeFails() {
	orig := httpListenAndServe
	defer func() {
		httpListenAndServe = orig
	}()
	httpListenAndServe = func(addr string, handler http.Handler) error {
		return fmt.Errorf("This is an error")
	}

	actual := New("myIp", "1234").Execute()

	s.Error(actual)
}

// ServeHTTP

func (s *ServerTestSuite) Test_ServeHTTP_ReturnsStatus404WhenURLIsUnknown() {
	req, _ := http.NewRequest("GET", "/this/url/does/not/exist", nil)

	srv := Serve{}
	srv.ServeHTTP(s.ResponseWriter, req)

	s.ResponseWriter.AssertCalled(s.T(), "WriteHeader", 404)
}

func (s *ServerTestSuite) Test_ServeHTTP_SetsContentTypeToJSON_WhenUrlIsJobAndMethodIsPut() {
	var actual string
	httpWriterSetContentType = func(w http.ResponseWriter, value string) {
		actual = value
	}
	req, _ := http.NewRequest("PUT", "/v1/docker-flow-cron/job", nil)
	cMock := CronerMock{
		AddJobMock: func(data cron.JobData) error {
			return nil
		},
	}

	srv := Serve{Cron: cMock}
	srv.ServeHTTP(s.ResponseWriter, req)

	s.Equal("application/json", actual)
}

func (s *ServerTestSuite) Test_ServeHTTP_InvokesCronAddJob() {
	expectedData := cron.JobData{
		Name: "my-job",
		Image: "my-image",
		Schedule: "@yearly",
	}
	actualData := cron.JobData{}
	js, _ := json.Marshal(expectedData)
	req, _ := http.NewRequest(
		"PUT",
		"/v1/docker-flow-cron/job",
		strings.NewReader(string(js)),
	)
	cMock := CronerMock{
		AddJobMock: func(data cron.JobData) error {
			actualData = data
			return nil
		},
	}

	srv := Serve{Cron: cMock}
	srv.ServeHTTP(s.ResponseWriter, req)

	s.Equal(expectedData, actualData)
}

func (s *ServerTestSuite) Test_ServeHTTP_ReturnsBadRequestWhenBodyIsNil() {
	req, _ := http.NewRequest("PUT", "/v1/docker-flow-cron/job", nil)
	cMock := CronerMock{
		AddJobMock: func(data cron.JobData) error {
			return nil
		},
	}

	srv := Serve{Cron: cMock}
	srv.ServeHTTP(s.ResponseWriter, req)

	s.ResponseWriter.AssertCalled(s.T(), "WriteHeader", 400)
}

func (s *ServerTestSuite) Test_ServeHTTP_InvokesInternalServerError_WhenAddJobFails() {
	expectedData := cron.JobData{
		Name: "my-job",
		Image: "my-image",
		Schedule: "@yearly",
	}
	js, _ := json.Marshal(expectedData)
	req, _ := http.NewRequest(
		"PUT",
		"/v1/docker-flow-cron/job",
		strings.NewReader(string(js)),
	)
	cMock := CronerMock{
		AddJobMock: func(data cron.JobData) error {
			return fmt.Errorf("This is an error")
		},
	}

	srv := Serve{Cron: cMock}
	srv.ServeHTTP(s.ResponseWriter, req)

	s.ResponseWriter.AssertCalled(s.T(), "WriteHeader", 500)
}

// Mock

type ResponseWriterMock struct {
	mock.Mock
}

func (m *ResponseWriterMock) Header() http.Header {
	m.Called()
	return make(map[string][]string)
}

func (m *ResponseWriterMock) Write(data []byte) (int, error) {
	params := m.Called(data)
	return params.Int(0), params.Error(1)
}

func (m *ResponseWriterMock) WriteHeader(header int) {
	m.Called(header)
}

func getResponseWriterMock() *ResponseWriterMock {
	mockObj := new(ResponseWriterMock)
	mockObj.On("Header").Return(nil)
	mockObj.On("Write", mock.Anything).Return(0, nil)
	mockObj.On("WriteHeader", mock.Anything)
	return mockObj
}

type CronerMock struct {
	AddJobMock func(data cron.JobData) error
	StopMock   func()
}

func (m CronerMock) AddJob(data cron.JobData) error {
	return m.AddJobMock(data)
}

func (m CronerMock) Stop() {
	m.StopMock()
}