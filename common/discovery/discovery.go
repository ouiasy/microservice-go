package discovery

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	consul "github.com/hashicorp/consul/api"
)

// Registry defines a service registry.
type Registry interface {
	// Register creates a service instance record in the registry.
	Register(ctx context.Context, instanceID string, serviceName string, hostPort string) error
	// Deregister removes a service insttance record from the registry.
	Deregister(ctx context.Context, instanceID string, serviceName string) error
	// ServiceAddresses returns the list of addresses of active instances of the given service.
	ServiceAddresses(ctx context.Context, serviceID string) ([]string, error)
	// ReportHealthyState is a push mechanism for reporting healthy state to the registry.
	ReportHealthyState(instanceID string, serviceName string) error
}

var (
	// ErrNotFound is returned when no service addresses are found.
	ErrNotFound = errors.New("no service addresses found")
)

// ConsulClient defines a struct contains a consul client
type ConsulClient struct {
	client *consul.Client
}

// NewConsulClient creates a new Consul-Client
func NewConsulClient(addr string) (*ConsulClient, error) {
	config := &consul.Config{
		Address: addr,
	}
	client, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &ConsulClient{client: client}, nil
}

// Register creates a service record in the registry.
func (c *ConsulClient) Register(ctx context.Context, instanceID string, serviceName string, hostPort string) error {
	parts := strings.Split(hostPort, ":")
	if len(parts) != 2 {
		return errors.New("hostPort must be in a form of <host>:<port>, example: localhost:8081")
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}
	return c.client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		Address: parts[0],
		ID:      instanceID,
		Name:    serviceName,
		Port:    port,
		Check:   &consul.AgentServiceCheck{CheckID: instanceID, TTL: "5s"},
	})
}

// Deregister removes a service record from the registry.
func (c *ConsulClient) Deregister(ctx context.Context, instanceID string, _ string) error {
	return c.client.Agent().ServiceDeregister(instanceID)
}

// ServiceAddresses returns the list of addresses of active instances of the given service.
func (c *ConsulClient) ServiceAddresses(ctx context.Context, serviceName string) ([]string, error) {
	entries, _, err := c.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	} else if len(entries) == 0 {
		return nil, ErrNotFound
	}
	var res []string
	for _, e := range entries {
		res = append(res, fmt.Sprintf("%s:%d", e.Service.Address, e.Service.Port))
	}
	return res, nil
}

// ReportHealthyState is a push mechanism for reporting healthy state to the registry.
func (c *ConsulClient) ReportHealthyState(instanceID string, _ string) error {
	return c.client.Agent().PassTTL(instanceID, "")
}

// GenerateInstanceID generates a pseudo-unique service instance identifier, using a service name
// suffixed by dash and a random number.
func GenerateInstanceID(serviceName string) string {
	return fmt.Sprintf("%s-%d", serviceName, rand.New(rand.NewSource(time.Now().UnixNano())).Int())
}
