#!/bin/bash
set -e
kubectl delete deployment manager --ignore-not-found
kubectl delete service manager --ignore-not-found
kubectl delete pod -l app=manager --force --grace-period=0 --ignore-not-found

kubectl delete deployment worker --ignore-not-found
kubectl delete pod -l app=worker --force --grace-period=0 --ignore-not-found

kubectl delete statefulset rabbitmq --ignore-not-found
kubectl delete service rabbitmq --ignore-not-found
kubectl delete service rabbitmq-headless --ignore-not-found
kubectl delete pod -l app=rabbitmq --force --grace-period=0 --ignore-not-found
kubectl delete pvc -l app=rabbitmq

kubectl delete statefulset mongodb --ignore-not-found
kubectl delete service mongodb-headless --ignore-not-found
kubectl delete pod -l app=mongodb --force --grace-period=0 --ignore-not-found
kubectl delete pvc -l app=mongodb

kubectl delete pv --all