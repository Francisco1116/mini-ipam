# Mini IPAM (IP Address Management)

A lightweight, cloud-native IP Address Management microservice built with Go. 
Designed to simulate core infrastructure orchestration by dynamically allocating non-overlapping subnets from a parent Virtual Private Cloud (VPC).

## 🚀 Why This Project?
In large-scale cloud environments (like OpenStack or Kubernetes), automated network allocation is critical. Manual IP management leads to conflicts and system downtime. This service mimics the backend orchestration logic that automatically calculates and issues safe, isolated CIDR blocks for new instances or container networks.

## 🛠 Tech Stack
* **Language:** Go (Golang) - utilizing `net/netip` for low-level bitwise network calculations.
* **Framework:** Gin (RESTful API routing)
* **Infrastructure:** Docker (Multi-stage build), Kubernetes (Deployment & LoadBalancer Service)
* **Architecture:** Cloud-Native, Stateless Microservice

## ✨ Core Features
* **Contract-First API:** RESTful endpoints designed with clear request/response schemas.
* **CIDR Collision Detection:** Linear search algorithm with bitwise IP-to-integer conversion to prevent overlapping subnets.
* **High Availability:** Deployed via Kubernetes with a replica set to ensure self-healing and load balancing.
* **Optimized Container:** Docker image size reduced to <20MB using multi-stage compilation (`scratch`/`alpine`).

## 📦 Quick Start (Kubernetes)

1. **Clone the repository:**
```bash
   git clone [https://github.com/Francisco1116/mini-ipam.git](https://github.com/Francisco1116/mini-ipam.git)
   cd mini-ipam
```

2. **Deploy to local K8s cluster:**
```bash
   kubectl apply -f k8s.yaml
```

3. **Verify the pods are running:**
```bash
   kubectl get pods
```

🔗 API Endpoints

1. **Create a VPC (Network)**
Creates a large parent network block.

POST /api/v1/networks
Body: {"name": "production-vpc", "cidr": "10.0.0.0/16"}

2. **Allocate a Subnet**
Dynamically calculates and returns an available subnet of the requested size.

POST /api/v1/networks/{id}/subnets
Body: {"name": "db-subnet", "prefix": 24}

3. **Get Network Status**
Returns the VPC details and all allocated subnets.

GET /api/v1/networks/{id}
Built to demonstrate foundational cloud infrastructure and networking concepts.