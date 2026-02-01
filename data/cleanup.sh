#!/bin/bash

echo "Starting MongoDB cleanup from Kubernetes..."
echo "Do you want to proceed? (y/n)"
read proceed
if [[ "$proceed" != "y" ]]; then
  echo "Cleanup aborted."
  exit 1
fi
# Delete MongoDB resources from Kubernetes
kubectl delete -f mongo/k8s/60-seed-job.yaml
kubectl delete -f mongo/k8s/50-seed-configmap.yaml
kubectl delete -f mongo/k8s/40-service-nodeport.yaml
kubectl delete -f mongo/k8s/30-service.yaml
kubectl delete -f mongo/k8s/20-deployment.yaml
kubectl delete -f mongo/k8s/10-pvc.yaml
kubectl delete -f mongo/k8s/00-namespace.yaml

echo "MongoDB cleanup completed."