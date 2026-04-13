#!/bin/bash
set -e
# # Включаем провиженер и дефолтный класс хранения
# minikube addons enable storage-provisioner
# minikube addons enable default-storageclass

# # Перезапускаем провиженер (на случай, если он "залип")
# minikube addons disable storage-provisioner
# minikube addons enable storage-provisioner

echo "🚀 Deploying CrackHash to Kubernetes..."

if ! kubectl get secret app-secrets  >/dev/null 2>&1; then
  echo "🔐 Creating secrets..."
  bash create-secrets.sh crackhash
fi

kubectl apply -f configmap.yaml

echo "🗄️  Deploying MongoDB..."
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

echo "🐰 Deploying RabbitMQ..."
kubectl apply -f rabbitmq/
kubectl rollout status statefulset/rabbitmq  --timeout=120s

echo "⚙️  Deploying manager..."
kubectl apply -f manager/
kubectl rollout status deployment/manager  --timeout=60s

echo "🔧 Deploying workers..."
kubectl apply -f worker/
kubectl rollout status deployment/worker  --timeout=60s

echo "✅ Deployment complete!"
echo "🔍 Check status: kubectl get pods "
echo "🌐 Access manager: kubectl port-forward svc/manager 8080:80 "