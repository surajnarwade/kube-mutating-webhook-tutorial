### HOW TO TEST ?


* deploy webhook

* once successfully deployed, deploy sample pod

```
kubectl run alpine --image=alpine --restart=Never  --overrides='{"apiVersion":"v1","metadata":{"annotations":{"sidecar-injector-webhook.morven.me/inject":"yes"}}}' --command -- sleep infinity
```

* it should successfully add sidecar to pod

```
 $ kubectl get pods
NAME      READY   STATUS              RESTARTS   AGE
alpine    0/2     ContainerCreating   0          55s
```

* manually renew the certs

```
 kubectl cert-manager renew -n sidecar-injector sidecar-injector
```

* redeploy pod

```
kubectl run alpine2 --image=alpine --restart=Never  --overrides='{"apiVersion":"v1","metadata":{"annotations":{"sidecar-injector-webhook.morven.me/inject":"yes"}}}' --command -- sleep infinity
```

* it shouldn't add sidecar

```
$ kubectl get pods
NAME      READY   STATUS              RESTARTS   AGE
alpine    0/2     ContainerCreating   0          55s
alpine1   1/1     Running             0          26s
```
and you should see error in webhook logs,

```
2020/07/13 15:21:24 http: TLS handshake error from 192.168.122.211:35086: remote error: tls: bad certificate
```

* approx after 30-40 seconds, you should see in webhook logs

```
2020/07/13 15:21:36 File Changed, reloading TLS certificate and key
```

* certs are renewed and now it should inject sidecar

```
$ kubectl run alpine3 --image=alpine --restart=Never  --overrides='{"apiVersion":"v1","metadata":{"annotations":{"sidecar-injector-webhook.morven.me/inject":"yes"}}}' --command -- sleep infinity
```

* wohooo


```
$ kubectl get pods
NAME      READY   STATUS              RESTARTS   AGE
alpine    0/2     ContainerCreating   0          55s
alpine1   1/1     Running             0          26s
alpine3   0/2     ContainerCreating   0          4s
```
