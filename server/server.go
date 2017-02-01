package server

import (
	"../cron"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"log"
)

type Serve struct {
	IP   string
	Port string
	Cron cron.Croner
}

type Response struct {
	Status  string
	Message string
}

var httpListenAndServe = http.ListenAndServe
var httpWriterSetContentType = func(w http.ResponseWriter, value string) {
	w.Header().Set("Content-Type", value)
}

var New = func(ip, port string) *Serve {
	return &Serve{
		IP:   ip,
		Port: port,
		Cron: cron.New(),
	}
}

func (s *Serve) Execute() error {
	log.Printf("Starting Web server running on %s:%s\n", s.IP, s.Port)
	address := fmt.Sprintf("%s:%s", s.IP, s.Port)
	if err := httpListenAndServe(address, s); err != nil {
		return err
	}
	return nil
}

func (s *Serve) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/v1/docker-flow-cron/job":
		s.addCronJob(w, req)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (s *Serve) addCronJob(w http.ResponseWriter, req *http.Request) {
	response := Response{
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
