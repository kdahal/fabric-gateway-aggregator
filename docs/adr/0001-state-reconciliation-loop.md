# ADR 001: Implementation of a State Reconciliation Loop for Fabric Port Sync

## Status
**Proposed** (March 16, 2026)

## Context
The **FP-Adapt-Aggregator** is responsible for maintaining the synchronization between the real-time physical status of 10,000+ fiber ports and our **Digital Twin** (hosted on GCP Spanner). 

Currently, the system relies on "push" updates (telemetry) from hardware. However, these are inherently unreliable due to:
* **UDP/Packet Loss:** Network congestion can cause status updates to be dropped.
* **Race Conditions:** Rapid state changes may arrive out of order.
* **Cold Starts:** If the Aggregator restarts, it loses the current "in-memory" state of the physical network.

Without a secondary verification mechanism, the Digital Twin and the Physical Port frequently fall out of sync, leading to "ghost ports" in the **FP-Remote-GUI**.

## Decision
We will implement a **State Reconciliation Loop** developed in **Go**.

### Mechanism
* **Worker Pool Pattern:** A dedicated Go routine pool will manage "polling" tasks to avoid exhausting system resources.
* **Low-Frequency Delta Checks:** The loop will poll device states at 1-minute intervals.
* **Comparison Logic:** The system will fetch the "Intended State" from GCP Spanner and compare it against the "Actual State" retrieved via Southbound APIs (gNMI/NETCONF).

### Healing Strategy
* If a mismatch is detected (e.g., Database = `UP`, Hardware = `DOWN`), the Aggregator will:
    1. Update the Digital Twin in GCP Spanner to reflect reality.
    2. Emit a high-priority event to the **FP-Remote-GUI** via the streaming layer (Pulsar/Kafka).
    3. Log the discrepancy for trend analysis (identifying faulty hardware).

## Considered Options

### 1. Pure Event-Driven Architecture (Rejected)
* **Reason:** While highly scalable, it lacks a "source of truth" verification. If a single Pulsar message is missed during a network spike, the UI remains permanently incorrect until the next state change.

### 2. Manual Synchronization (Rejected)
* **Reason:** Requires human intervention or a CLI trigger. This is incompatible with Lumen’s "Network-as-a-Service" automation goals and creates a bottleneck for operational teams.

## Consequences

### Pros
* **Guaranteed Eventual Consistency:** The system will always self-correct within 60 seconds.
* **Reliability:** The UI becomes a "Trusted Source" for technicians.
* **Observability:** We gain data on how often hardware state drifts from software intent.

### Cons
* **Increased Southbound Traffic:** Frequent polling adds load to the switches. 
* **Mitigation:** We will use Go’s `time.Ticker` and throttled concurrency to ensure we do not overwhelm older legacy hardware.

---
**Approvers:** * Kumar Dahal (Lead Architect)