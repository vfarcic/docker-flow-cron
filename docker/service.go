package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/swarm"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"fmt"
)

const dockerApiVersion = "v1.24"

type Servicer interface {
	GetServices(jobName string) ([]swarm.Service, error)
	GetTasks(jobName string) ([]swarm.Task, error)
}

type Service struct {
	Client *client.Client
}

func New(host string) (*Service, error) {
	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	c, err := client.NewClient(host, dockerApiVersion, nil, defaultHeaders)
	if err != nil {
		return &Service{}, err
	}
	return &Service{Client: c}, nil
}

func (s *Service) GetServices(jobName string) ([]swarm.Service, error) {
	filter := filters.NewArgs()
	filter.Add("label", "com.df.cron=true")
	if len(jobName) > 0 {
		filter.Add("label", fmt.Sprintf("com.df.cron.name=%s", jobName))
	}
	services, err := s.Client.ServiceList(context.Background(), types.ServiceListOptions{Filters: filter})
	if err != nil {
		return []swarm.Service{}, err
	}
	return services, nil
}

func (s *Service) GetTasks(jobName string) ([]swarm.Task, error) {
	filter := filters.NewArgs()
	filter.Add("label", "com.df.cron=true")
	filter.Add("label", fmt.Sprintf("com.df.cron.name=%s", jobName))
	tasks, err := s.Client.TaskList(context.Background(), types.TaskListOptions{Filters: filter})
	if err != nil {
		return []swarm.Task{}, err
	}
	return tasks, nil
}
