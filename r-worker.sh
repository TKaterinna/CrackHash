#!/bin/bash
set -e

kubectl delete deployment worker --ignore-not-found
kubectl delete pod -l app=worker --force --grace-period=0 --ignore-not-found

kubectl apply -f worker/
kubectl rollout status deployment/worker  --timeout=60s

kubectl get pod -l app=worker