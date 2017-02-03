package server

import (
	"../cron"
	"../docker"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"github.com/gorilla/mux"
)

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

type ResponseDetails struct {
	Status  string
	Message string
	Job     cron.JobData
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
	return &Serve{
		IP:      ip,
		Port:    port,
		Cron:    cron.New(),
		Service: service,
	}, nil
}

func (s *Serve) Execute() error {
	log.Printf("Starting Web server running on %s:%s\n", s.IP, s.Port)
	address := fmt.Sprintf("%s:%s", s.IP, s.Port)
	// TODO: Test routes
	r := mux.NewRouter()
	sr := r.PathPrefix("/v1/docker-flow-cron/job").Subrouter()
	sr.HandleFunc("/", s.JobPutHandler).Methods("PUT")
	sr.HandleFunc("/", s.JobGetHandler).Methods("GET")
//	sr.HandleFunc("/{name}", s.JobDetailsHandler).Methods("GET")
	if err := httpListenAndServe(address, r); err != nil {
		return err
	}
	return nil
}

// TODO: Develop
//func (s *Serve) JobDetailsHandler(w http.ResponseWriter, req *http.Request) {
////	vars := mux.Vars(req)
////	println(vars["name"])
//	response := ResponseDetails{
//		Status: "OK",
//	}
//	httpWriterSetContentType(w, "application/json")
//	js, _ := json.Marshal(response)
//	w.Write(js)
//}

func (s *Serve) JobGetHandler(w http.ResponseWriter, req *http.Request) {
	services, _ := s.Service.GetServices()
	jobs := map[string]cron.JobData{}
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
		jobs[name] = cron.JobData{
			Name:     name,
			Image:    service.Spec.TaskTemplate.ContainerSpec.Image,
			Command:  service.Spec.Annotations.Labels["com.df.cron.command"],
			Schedule: service.Spec.Annotations.Labels["com.df.cron.schedule"],
		}
	}
	response := Response{
		Status: "OK",
		Jobs:   jobs,
	}
	httpWriterSetContentType(w, "application/json")
	js, _ := json.Marshal(response)
	w.Write(js)
}

func (s *Serve) JobPutHandler(w http.ResponseWriter, req *http.Request) {
	// TODO: Add job to the response
	response := ResponseDetails{
		Status: "OK",
	}
	if req.Body == nil {
		w.WriteHeader(http.StatusBadRequest)
		response.Message = "Request body is mandatory"
		response.Status = "NOK"
	} else {
		defer func() { req.Body.Close() }()
		body, _ := ioutil.ReadAll(req.Body)
		data := cron.JobData{}
		json.Unmarshal(body, &data)
		if err := s.Cron.AddJob(data); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			response.Message = err.Error()
			response.Status = "NOK"
		}
	}
	httpWriterSetContentType(w, "application/json")
	js, _ := json.Marshal(response)
	w.Write(js)
}
