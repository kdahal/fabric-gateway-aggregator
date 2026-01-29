resource "google_compute_network" "fabric_vpc" {
  name                    = "fabric-access-ecosystem"
  auto_create_subnetworks = false
}

# Example of a custom provider block for the Fabric Controller
resource "fabric_port_config" "edge_port_01" {
  port_id = "10.0.0.1/port-1"
  vlan    = 100
}