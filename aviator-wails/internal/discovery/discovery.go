package discovery

import (
	"log"
	"os"


	"github.com/grandcat/zeroconf"
)

type DiscoveryService struct {
	server *zeroconf.Server
}

func NewDiscoveryService(port int) (*DiscoveryService, error) {
	hostname, _ := os.Hostname()
	// Service name: "Aviator At [Hostname]"
	serviceName := "Aviator At " + hostname
	serviceType := "_aviator._tcp"
	domain := "local."

	// Metadata to help client if needed
	txtRecords := []string{"version=1.0", "type=aviator-go"}

	// Register the service
	// We use "nil" for interface to bind to all available interfaces
	server, err := zeroconf.Register(serviceName, serviceType, domain, port, txtRecords, nil)
	if err != nil {
		return nil, err
	}

	log.Printf("[Discovery] Registered mDNS service: %s (%s) on port %d", serviceName, serviceType, port)

	return &DiscoveryService{
		server: server,
	}, nil
}

func (ds *DiscoveryService) Shutdown() {
	if ds.server != nil {
		ds.server.Shutdown()
	}
}
