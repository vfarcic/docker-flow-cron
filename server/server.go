package server

import (
	"../cron"
	"../docker"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types/swarm"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var muxVars = mux.Vars

type Serve struct {
	IP      string
	Port    string
	Cron    cron.Croner
	Service docker.Servicer
}

type Response struct {
	Status  string
	Message string
	Jobs    map[string]cron.JobData
}

type Execution struct {
	ServiceId string
	CreatedAt time.Time
	Status    swarm.TaskStatus
}

type ResponseDetails struct {
	Status     string
	Message    string
	Job        cron.JobData
	Executions []Execution
}

var httpListenAndServe = http.ListenAndServe
var httpWriterSetContentType = func(w http.ResponseWriter, value string) {
	w.Header().Set("Content-Type", value)
}

var New = func(ip, port, dockerHost string) (*Serve, error) {
	service, err := docker.New(dockerHost)
	if err != nil {
		return &Serve{}, err
	}
	cron, _ := cron.New(dockerHost)
	return &Serve{
		IP:      ip,
		Port:    port,
		Cron:    cron,
		Service: service,
	}, nil
}

func (s *Serve) Execute() error {
	fmt.Printf("Starting Web server running on %s:%s\n", s.IP, s.Port)
	address := fmt.Sprintf("%s:%s", s.IP, s.Port)
	// TODO: Test routes
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/v1/docker-flow-cron/job", s.JobGetHandler).Methods("GET")
	r.HandleFunc("/v1/docker-flow-cron/job/{jobName}", s.JobPutHandler).Methods("PUT")
	r.HandleFunc("/v1/docker-flow-cron/job/{jobName}", s.JobDetailsHandler).Methods("GET")
	// TODO: Document
	r.HandleFunc("/v1/docker-flow-cron/job/{jobName}", s.JobDeleteHandler).Methods("DELETE")
	if err := httpListenAndServe(address, r); err != nil {
		return err
	}
	return nil
}

func (s *Serve) JobDeleteHandler(w http.ResponseWriter, req *http.Request) {
	jobName := muxVars(req)["jobName"]
	response := ResponseDetails{
		Status:  "OK",
		Message: fmt.Sprintf("%s was deleted", jobName),
	}
	err := s.Cron.RemoveJob(jobName)
	if err != nil {
		response.Status = "NOK"
		response.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
	js, _ := json.Marshal(response)
	w.Write(js)
}

func (s *Serve) JobDetailsHandler(w http.ResponseWriter, req *http.Request) {
	jobName := muxVars(req)["jobName"]
	services, err := s.Service.GetServices(jobName)
	httpWriterSetContentType(w, "application/json")
	response := ResponseDetails{
		Status:     "OK",
		Message:    "",
		Job:        cron.JobData{},
		Executions: []Execution{},
	}
	if err != nil {
		response.Status = "NOK"
		response.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		index := len(services) - 1
		if index < 0 {
			response.Status = "NOK"
			response.Message = "Could not find the job"
			w.WriteHeader(http.StatusNotFound)
		} else {
			service := services[index]
			tasks, err := s.Service.GetTasks(jobName)
			if err != nil {
				response.Status = "NOK"
				response.Message = err.Error()
				w.WriteHeader(http.StatusInternalServerError)
			} else {
				executions := []Execution{}
				for _, t := range tasks {
					execution := Execution{
						ServiceId: t.ServiceID,
						CreatedAt: t.CreatedAt,
						Status:    t.Status,
					}
					executions = append(executions, execution)
				}
				response.Job = s.getJob(service)
				response.Executions = executions
			}
		}
	}
	js, _ := json.Marshal(response)
	w.Write(js)
}

func (s *Serve) JobGetHandler(w http.ResponseWriter, req *http.Request) {
	response := Response{
		Status: "OK",
	}
	jobs, err := s.Cron.GetJobs()
	if err != nil {
		response.Status = "NOK"
		response.Message = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
	}
	response.Jobs = jobs
	httpWriterSetContentType(w, "application/json")
	js, _ := json.Marshal(response)
	w.Write(js)
}

func (s *Serve) JobPutHandler(w http.ResponseWriter, req *http.Request) {
	jobName := muxVars(req)["jobName"]
	response := ResponseDetails{
		Status: "OK",
	}
	if req.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Status = "NOK"
		response.Message = "Request body is mandatory"
	} else {
		defer func() { req.Body.Close() }()
		body, _ := ioutil.ReadAll(req.Body)
		data := cron.JobData{}
		json.Unmarshal(body, &data)
		data.Name = jobName
		response.Job = data
		if err := s.Cron.AddJob(data); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response.Status = "NOK"
			response.Message = err.Error()
		} else {
			response.Message = fmt.Sprintf("Job %s has been scheduled", data.Name)
		}
	}
	httpWriterSetContentType(w, "application/json")
	js, _ := json.Marshal(response)
	w.Write(js)
}

// TODO: Remove when JobDetailsHandler is refactored to use cron
func (s *Serve) getJob(service swarm.Service) cron.JobData {
	command := ""
	for _, v := range service.Spec.TaskTemplate.ContainerSpec.Args {
		if strings.Contains(v, " ") {
			command = fmt.Sprintf(`%s "%s"`, command, v)
		} else {
			command = fmt.Sprintf(`%s %s`, command, v)
		}
	}
	name := service.Spec.Annotations.Labels["com.df.cron.name"]
	return cron.JobData{
		Name:     name,
		ServiceName: service.Spec.Annotations.Labels["com.df.cron.servicename"],
		Image:    service.Spec.TaskTemplate.ContainerSpec.Image,
		Command:  service.Spec.Annotations.Labels["com.df.cron.command"],
		Schedule: service.Spec.Annotations.Labels["com.df.cron.schedule"],
	}
}
