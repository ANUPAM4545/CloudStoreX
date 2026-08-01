# CloudStoreX GitOps Application Configurations (ArgoCD)

This directory contains ArgoCD Application manifests for deploying CloudStoreX using GitOps practices across multiple environments (`development` and `production`).

## Directory Structure
- `cloudstorex-dev-app.yaml`: Application manifest pointing to `deploy/kubernetes/overlays/development/`.
- `cloudstorex-prod-app.yaml`: Application manifest pointing to `deploy/kubernetes/overlays/production/`.

## Deploying with ArgoCD
```bash
kubectl apply -f argocd/cloudstorex-dev-app.yaml
kubectl apply -f argocd/cloudstorex-prod-app.yaml
```
