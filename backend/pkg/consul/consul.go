package consul

import (
	"fmt"
	"time"

	"github.com/hashicorp/consul/api"
)

type Client struct {
	*api.Client
}

func New(addr, token string) (*Client, error) {
	cfg := api.DefaultConfig()
	cfg.Address = addr
	cfg.Token = token

	cli, err := api.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("consul new: %w", err)
	}
	return &Client{cli}, nil
}

func (c *Client) RegisterService(name, id, host string, port int, tags []string) error {
	reg := &api.AgentServiceRegistration{
		ID:      id,
		Name:    name,
		Tags:    tags,
		Port:    port,
		Address: host,
		Check: &api.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("http://%s:%d/health", host, port),
			Interval:                       "10s",
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: "30s",
		},
	}
	return c.Client.Agent().ServiceRegister(reg)
}

func (c *Client) Deregister(id string) error {
	return c.Client.Agent().ServiceDeregister(id)
}

func (c *Client) GetService(name string) ([]*api.ServiceEntry, error) {
	services, _, err := c.Client.Health().Service(name, "", true, nil)
	return services, err
}

func (c *Client) WatchServices(name string, stopCh <-chan struct{}) (<-chan []*api.ServiceEntry, error) {
	ch := make(chan []*api.ServiceEntry)
	go func() {
		defer close(ch)
		opts := &api.QueryOptions{WaitIndex: 0, WaitTime: 10 * time.Second}
		for {
			select {
			case <-stopCh:
				return
			default:
				services, meta, err := c.Client.Health().Service(name, "", true, opts)
				if err != nil {
					time.Sleep(5 * time.Second)
					continue
				}
				opts.WaitIndex = meta.LastIndex
				ch <- services
			}
		}
	}()
	return ch, nil
}
