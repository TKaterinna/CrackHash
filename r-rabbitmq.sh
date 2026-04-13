#!/bin/bash
set -e

kubectl delete statefulset rabbitmq --ignore-not-found
kubectl delete service rabbitmq --ignore-not-found
kubectl delete service rabbitmq-headless --ignore-not-found
kubectl delete pod -l app=rabbitmq --force --grace-period=0 --ignore-not-found
kubectl delete pvc -l app=rabbitmq

kubectl apply -f rabbitmq/
kubectl rollout status statefulset/rabbitmq  --timeout=120s

kubectl get pod -l app=rabbitmq