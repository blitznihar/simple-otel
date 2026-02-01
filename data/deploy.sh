#!/bin/bash

echo "Starting MongoDB deployment to Kubernetes..."
echo "Do you want to proceed? (y/n)"
read proceed
if [[ "$proceed" != "y" ]]; then
  echo "Deployment aborted."
  exit 1
fi

# Deploy MongoDB to Kubernetes
kubectl apply -f mongo/k8s/00-namespace.yaml
kubectl apply -f mongo/k8s/10-pvc.yaml
kubectl apply -f mongo/k8s/20-deployment.yaml
kubectl apply -f mongo/k8s/30-service.yaml
kubectl apply -f mongo/k8s/40-service-nodeport.yaml

# Verify deployment
kubectl -n simple-otel get pods,svc,pvc
kubectl -n simple-otel logs deploy/mongo --tail=50


kubectl -n simple-otel run tmp --rm -it --image=mongo:7 --restart=Never -- \
  mongosh "mongodb://mongo:27017" --eval "db.runCommand({ ping: 1 })"

echo "MongoDB deployment completed."
echo "You can access MongoDB at mongodb://localhost:31017"