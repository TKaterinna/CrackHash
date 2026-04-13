#!/bin/bash
set -e

kubectl delete statefulset mongodb --ignore-not-found
kubectl delete service mongodb-headless --ignore-not-found
kubectl delete pod -l app=mongodb --force --grace-period=0 --ignore-not-found
kubectl delete pvc -l app=mongodb

kubectl apply -f mongodb/
for i in 0 1 2; do
    kubectl wait --for=condition=ready pod/mongodb-$i --timeout=300s
done

# ⚠️ Инициализация реплика-сета (однократно!)
if ! kubectl exec mongodb-0 -- mongosh --eval "rs.status()" >/dev/null 2>&1; then
  echo "🔄 Initializing MongoDB replica set..."
  kubectl exec mongodb-0 -- mongosh --eval '
    rs.initiate({
      _id: "rs0",
      members: [
        { _id: 0, host: "mongodb-0.mongodb-headless.default.svc.cluster.local:27017" },
        { _id: 1, host: "mongodb-1.mongodb-headless.default.svc.cluster.local:27017" },
        { _id: 2, host: "mongodb-2.mongodb-headless.default.svc.cluster.local:27017" }
      ]
    })
  '
fi

kubectl get pod -l app=mongodb