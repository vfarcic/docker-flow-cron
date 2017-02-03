package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"log"
)

const dockerApiVersion = "v1.24"

type Servicer interface {
	GetServices() ([]swarm.Service, error)
}

type Service struct {
	Client *client.Client
}

func New(host string) (*Service, error) {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	c, err := client.NewClient(host, dockerApiVersion, nil, defaultHeaders)
	if err != nil {
		log.Printf(err.Error())
		return &Service{}, err
	}
	return &Service{Client: c}, nil
}

func (s *Service) GetServices() ([]swarm.Service, error) {
	filter := filters.NewArgs()
	filter.Add("label", "com.df.cron=true")
	services, err := s.Client.ServiceList(context.Background(), types.ServiceListOptions{Filters: filter})
	if err != nil {
		log.Printf(err.Error())
		return []swarm.Service{}, err
	}
	return services, nil
}
