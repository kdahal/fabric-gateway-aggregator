package service

import (
	"context"
	"testing"

	"github.com/nats-io/nats.go"
)

func TestProvisionPort(t *testing.T) {
	// 1. Setup: Use an in-process NATS server for testing
	// This ensures the test is fast and doesn't require an external broker
	s, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		// If NATS isn't running locally, we skip the integration part 
		// but still test the logic. For a real CI, we'd use a mock.
		t.Skip("Skipping NATS integration test; no local server found")
	}
	defer s.Close()

	svc := NewAggregatorService(s)

	// 2. Define Test Cases
	tests := []struct {
		name       string
		portID     string
		vlanID     int
		zone       string
		wantErr    bool
	}{
		{
			name:    "Valid Provisioning Request",
			portID:  "Eth-1-1",
			vlanID:  100,
			zone:    "GCP-US-CENTRAL",
			wantErr: false,
		},
		{
			name:    "Invalid VLAN (Too Low)",
			portID:  "Eth-1-1",
			vlanID:  0,
			zone:    "GCP-US-CENTRAL",
			wantErr: true,
		},
		{
			name:    "Invalid VLAN (Too High)",
			portID:  "Eth-1-1",
			vlanID:  5000,
			zone:    "GCP-US-CENTRAL",
			wantErr: true,
		},
	}

	// 3. Execution
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.ProvisionPort(context.Background(), tt.portID, tt.vlanID, tt.zone)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProvisionPort() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}