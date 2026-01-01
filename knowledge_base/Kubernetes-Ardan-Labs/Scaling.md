#k8s

## Scaling DockerCoin App
### Adding more workers for more loops per second
```sh
$ kubectl scale deployment worker --replicas 2
```