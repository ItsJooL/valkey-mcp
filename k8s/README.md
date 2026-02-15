# Kubernetes Deployment Guide

Simple Kubernetes deployment for Valkey MCP Server.

## Prerequisites

- Kubernetes cluster running
- `kubectl` configured
- Valkey/Redis instance accessible in your cluster

## Quick Start

### 1. Deploy Everything

```bash
kubectl apply -f k8s/valkey-mcp-server.yaml
```

This creates:
- ConfigMap with Valkey connection settings
- Deployment with 1 replica
- ClusterIP Service (internal)
- NodePort Service (external access on port 30080)

### 2. Verify Deployment

```bash
# Check if pod is running
kubectl get pods -l app=valkey-mcp-server

# Check services
kubectl get svc -l app=valkey-mcp-server

# View logs
kubectl logs -l app=valkey-mcp-server --tail=50
```

### 3. Update Configuration

Edit the ConfigMap to change Valkey connection:

```bash
kubectl edit configmap valkey-mcp-config
```

Update values:
```yaml
data:
  VALKEY_URL: "valkey://your-valkey-host:6379"
  VALKEY_DB: "0"
  VALKEY_PASSWORD: "your-password"
```

### 4. Rollout Changes

After updating ConfigMap, restart the deployment:

```bash
kubectl rollout restart deployment/valkey-mcp-server
```

Or update the annotation to force a rollout:

```bash
kubectl patch deployment valkey-mcp-server \
  -p '{"spec":{"template":{"metadata":{"annotations":{"configmap/checksum":"'$(date +%s)'"}}}}}'
```

## Accessing the Server

### From Inside Cluster

```bash
# Using ClusterIP service
curl http://valkey-mcp-service:8080
```

### From Outside Cluster

```bash
# Using NodePort (port 30080)
curl http://<node-ip>:30080
```

## Scaling

```bash
# Scale to 3 replicas
kubectl scale deployment valkey-mcp-server --replicas=3

# Check status
kubectl get deployment valkey-mcp-server
```

## Troubleshooting

### Check Pod Status

```bash
kubectl get pods -l app=valkey-mcp-server
kubectl describe pod -l app=valkey-mcp-server
```

### View Logs

```bash
# Latest logs
kubectl logs -l app=valkey-mcp-server

# Follow logs
kubectl logs -l app=valkey-mcp-server -f

# Previous crashed container
kubectl logs -l app=valkey-mcp-server --previous
```

### Check ConfigMap

```bash
kubectl get configmap valkey-mcp-config -o yaml
```

### Restart Deployment

```bash
kubectl rollout restart deployment/valkey-mcp-server
kubectl rollout status deployment/valkey-mcp-server
```

## Clean Up

```bash
# Delete everything
kubectl delete -f k8s/valkey-mcp-server.yaml

# Or delete individually
kubectl delete deployment valkey-mcp-server
kubectl delete service valkey-mcp-service valkey-mcp-nodeport
kubectl delete configmap valkey-mcp-config
```

## Configuration Options

### ConfigMap Values

| Key | Description | Default |
|-----|-------------|---------|
| `VALKEY_URL` | Valkey connection URL | `valkey://valkey-service:6379` |
| `VALKEY_DB` | Database number (0-15) | `0` |
| `VALKEY_PASSWORD` | Authentication password | `` |

### Resource Limits

Adjust in the deployment spec:

```yaml
resources:
  requests:
    cpu: 100m      # Minimum CPU
    memory: 128Mi  # Minimum memory
  limits:
    cpu: 500m      # Maximum CPU
    memory: 512Mi  # Maximum memory
```

## Example: Complete Workflow

```bash
# 1. Deploy
kubectl apply -f k8s/valkey-mcp-server.yaml

# 2. Wait for pod to be ready
kubectl wait --for=condition=ready pod -l app=valkey-mcp-server --timeout=60s

# 3. Check status
kubectl get all -l app=valkey-mcp-server

# 4. Update config
kubectl edit configmap valkey-mcp-config

# 5. Rollout changes
kubectl rollout restart deployment/valkey-mcp-server

# 6. Watch rollout
kubectl rollout status deployment/valkey-mcp-server

# 7. Verify
kubectl logs -l app=valkey-mcp-server --tail=20
```

## Notes

- The deployment uses a simple readiness probe checking if the process is running
- NodePort 30080 is used for external access (can be changed in the manifest)
- The pod runs as non-root user (UID 1000) for security
- ConfigMap changes require manual rollout restart
- No HPA, affinity, or PodDisruptionBudget for simplicity
