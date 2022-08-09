# go-endpoint

RUN 

```bash
$ docker-compose -f .\docker-compose-dev.yml up -d
```
KUBERNETES
```bash
$ kubectl apply -f .\deployment.yml
```

Forward for local
```bash
$ kubectl port-forward service/go-endpoint-svc  5555:5555
$ kubectl port-forward service/mongo-express-svc 5444:5444
```