package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

// PortTask represents the message payload for the distributed fabric
type PortTask struct {
	ID         string    `json:"port_id"`
	VlanID     int       `json:"vlan_id"`
	FabricZone string    `json:"fabric_zone"`
	Action     string    `json:"action"`
	CreatedAt  time.Time `json:"created_at"`
}

// AggregatorService manages the coordination between the API and the Network Fabric
type AggregatorService struct {
	nc *nats.Conn
}

// NewAggregatorService initializes a new service with a NATS connection
func NewAggregatorService(nc *nats.Conn) *AggregatorService {
	return &AggregatorService{nc: nc}
}

// ProvisionPort processes the logical request and dispatches an asynchronous task to the physical fabric
func (s *AggregatorService) ProvisionPort(ctx context.Context, portID string, vlanID int, zone string) error {
	// 1. Logic Layer: Validate business requirements (e.g., Fabric Port EcoSystem rules)
	if vlanID < 1 || vlanID > 4094 {
		return fmt.Errorf("invalid vlan %d: must be within range 1-4094", vlanID)
	}

	task := PortTask{
		ID:         portID,
		VlanID:     vlanID,
		FabricZone: zone,
		Action:     "PROVISION_UP",
		CreatedAt:  time.Now(),
	}

	payload, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to encode port task: %w", err)
	}

	// 2. Asynchronous Programming: Dispatch to the message broker
	// This ensures the API remains responsive (202 Accepted) while the hardware configures
	subject := fmt.Sprintf("fabric.zone.%s.provision", zone)
	if err := s.nc.Publish(subject, payload); err != nil {
		return fmt.Errorf("failed to dispatch to fabric broker: %w", err)
	}

	// 3. Optional: Logic for logging or Digital Twin state update
	fmt.Printf("Successfully dispatched provisioning task for port %s in zone %s\n", portID, zone)

	return nil
}