# Docker Compose Guide

This guide explains how to run the Ride Booking application using Docker Compose.

## Prerequisites

- Docker Engine 20.10.0+
- Docker Compose V2+

## Quick Start

1. **Clone the repository** (if you haven't already)

   ```bash
   git clone <repository-url>
   cd ride-booking-go
   ```

2. **Create environment file**

   ```bash
   cp .env.docker-compose.example .env
   ```

3. **Start all services**

   ```bash
   docker-compose up -d
   ```

4. **View logs**

   ```bash
   docker-compose logs -f
   ```

5. **Stop all services**

   ```bash
   docker-compose down
   ```

## Services Overview

The application consists of the following services:

| Service | Port | Description |
|---------|------|-------------|
| **RabbitMQ** | 5672, 15672 | Message broker (Management UI on 15672) |
| **API Gateway** | 8081 | HTTP API Gateway |
| **Trip Service** | 9093 | gRPC service for trip management |
| **Driver Service** | 9092 | gRPC service for driver management |
| **Payment Service** | 9004 | gRPC service for payment processing |
| **Web Frontend** | 3000 | Next.js web application |

## Environment Variables

The `.env` file contains the following configuration:

```bash
# Environment
ENVIRONMENT=development

# RabbitMQ
RABBITMQ_USER=guest
RABBITMQ_PASS=guest

# Stripe Payment
STRIPE_SECRET_KEY="your-stripe-secret-key"
STRIPE_SUCCESS_URL=http://localhost:3000?payment=success
STRIPE_CANCEL_URL=http://localhost:3000?payment=cancel

# Web Frontend
NODE_ENV=development
NEXT_PUBLIC_API_URL=http://localhost:8081
NEXT_PUBLIC_WEBSOCKET_URL=ws://localhost:8081/ws
```

## Common Commands

### Build and Start Services

```bash
# Build images and start all services
docker-compose up --build

# Start in detached mode
docker-compose up -d

# Start specific service
docker-compose up api-gateway
```

### View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f api-gateway

# Last 100 lines
docker-compose logs --tail=100 -f
```

### Stop Services

```bash
# Stop all services
docker-compose down

# Stop and remove volumes
docker-compose down -v

# Stop specific service
docker-compose stop api-gateway
```

### Restart Services

```bash
# Restart all services
docker-compose restart

# Restart specific service
docker-compose restart api-gateway
```

### Rebuild Services

```bash
# Rebuild all services
docker-compose build

# Rebuild specific service
docker-compose build api-gateway

# Rebuild without cache
docker-compose build --no-cache
```

### Execute Commands in Containers

```bash
# Execute shell in a container
docker-compose exec api-gateway sh

# Execute specific command
docker-compose exec api-gateway ps aux
```

## Access Points

Once all services are running:

- **Web Application**: <http://localhost:3000>
- **API Gateway**: <http://localhost:8081>
- **RabbitMQ Management UI**: <http://localhost:15672> (guest/guest)

## Troubleshooting

### Services won't start

```bash
# Check logs
docker-compose logs

# Rebuild images
docker-compose build --no-cache

# Clean up and restart
docker-compose down -v
docker-compose up --build
```

### Port conflicts

If you have services running on the same ports:

```bash
# Modify ports in docker-compose.yml
# Change "8081:8081" to "8082:8081" for example
```

### RabbitMQ connection issues

```bash
# Check RabbitMQ health
docker-compose exec rabbitmq rabbitmq-diagnostics check_running

# Restart RabbitMQ
docker-compose restart rabbitmq
```

### View service status

```bash
docker-compose ps
```

### Inspect networks

```bash
docker network ls
docker network inspect ride-booking-go_ride-booking-network
```

## Differences from Tilt

The Docker Compose setup differs from Tilt in the following ways:

### Tilt Setup

- Uses Kubernetes (minikube/kind) for orchestration
- Hot reloading with live_update for development
- Compiles Go binaries locally, copies to container
- Kubernetes resources (ConfigMaps, Secrets, Services)
- Port-forwarding through kubectl

### Docker Compose Setup

- Uses Docker networking for service communication
- Multi-stage Docker builds (compile inside container)
- Simpler networking model
- Environment variables instead of ConfigMaps/Secrets
- Direct port mapping to host

## Production Considerations

For production deployment:

1. **Update Stripe credentials** in `.env`
2. **Change RabbitMQ credentials** from default guest/guest
3. **Use production-grade persistence** for RabbitMQ
4. **Add health checks** and restart policies
5. **Configure resource limits** appropriately
6. **Use secrets management** instead of .env files
7. **Enable TLS/SSL** for external endpoints
8. **Set up proper logging** and monitoring

## Development Workflow

For active development, you may prefer Tilt for hot-reloading. Use Docker Compose for:

- Quick testing of the full stack
- CI/CD pipelines
- Integration testing
- Demonstration purposes
- Production-like local environment

## Clean Up

To completely remove all containers, networks, and volumes:

```bash
docker-compose down -v --remove-orphans
docker system prune -a
```
