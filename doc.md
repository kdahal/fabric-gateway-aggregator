
## High-Level Architecture
## 1. The Architectural Design: "Fabric-Port-as-a-Service"
Below is a high-level design using the Hexagonal Architecture (Ports and Adapters) pattern. This is the gold standard for Go and Node.js because it separates the business logic from the messy physical network.

# High-Level Component Flow
1. UI Layer (FP-Remote-GUI): A React/Node.js application where technicians or customers view port status.
2. Orchestration Layer (GCP): Where your logic lives. It validates requests (e.g., "Does this user have permission to increase bandwidth?").
3. Aggregation Layer (FP-Adapt-Aggregator): This is the critical engine. It "aggregates" thousands of raw status updates from physical switches and "adapts" them into a clean API format.

```mermaid
graph TD
    subgraph "External World"
        User((Admin/Customer))
    end

    subgraph "Lumen Ecosystem (GCP)"
        GUI[FP-Remote-GUI / Node.js]
        Agg[FP-Adapt-Aggregator / Go]
        Twin[(Digital Twin / Spanner Graph)]
    end

    subgraph "Physical Network (On-Prem/Edge)"
        Switch1[Cisco/Arista Switch]
        Switch2[Juniper Router]
        Sensors[Telemetry Stream]
    end

    User -->|OpenAPI| GUI
    GUI -->|gRPC/REST| Agg
    Agg <--> Twin
    Agg -->|Pulsar/Kafka| Sensors
    Agg -->|Netconf/YANG| Switch1
    Agg -->|Netconf/YANG| Switch2

```