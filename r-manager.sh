#!/bin/bash
set -e

kubectl delete deployment manager --ignore-not-found
kubectl delete service manager --ignore-not-found
kubectl delete pod -l app=manager --force --grace-period=0 --ignore-not-found

kubectl apply -f manager/
kubectl rollout status deployment/manager  --timeout=60s

kubectl get pod -l app=manager