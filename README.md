# CrackHash
## Команды
- Прокидывание порта для реббита
kubectl port-forward svc/rabbitmq 15672:15672

- Прокидывание порта для графаны
kubectl port-forward -n monitoring svc/minimal-grafana 3000:80

- Прокидывание порта для постмана
kubectl port-forward svc/manager 8080:80

- Поднятие миникуба
minikube start
minikube start --memory=2000 --cpus=2 --driver=docker --alsologtostderr -v=2

- Проверка кубера
minikube dashboard

- Проверка поднятых подов
kubectl get pod

- Перезапуск мониторинга и его проверка
kubectl apply -f monitoring/minimal-stack.yaml
kubectl rollout restart deployment minimal-prometheus -n monitoring
kubectl get pods -n monitoring
kubectl port-forward -n monitoring svc/minimal-grafana 3000:80

- Статусы монго
kubectl exec mongodb-0 -- mongosh --quiet --eval '
  const status = rs.status();
  print("Replica Set: " + status.set + "\n");
  status.members.forEach(m => {
    const role = m.stateStr.padEnd(10);
    const health = m.health === 1 ? "🟢 UP" : "🔴 DOWN";
    const lag = m.optimeDate ? "" : " (no optime)";
    print("  " + m.name + " → " + role + " " + health + lag);
  });
'
