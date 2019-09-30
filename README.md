podlauncher
---

Quick tool written for a customer to launch similar pods quickly using K8s api that will
have their own K8s services.

Flags
---
* --prefix=[prefix of the pod/container names]
* --num=[number of containers to create]
* --image=[path to docker image to use]
* --namespace=[namespace to use, default is default]
* --port[port to use for the service]

Example
---
```podlauncher --prefix=nginx --num=3 --image=nginx --namespace=default --port=80```