podlauncher
---

Quick tool written for a customer to launch similar pods quickly using K8s api that will
have their own dedicated K8s services and possibly multiple ports.

Flags
---
* --prefix=[prefix of the pod/container names]
* --num=[number of containers to create]
* --image=[path to docker image to use]
* --namespace=[namespace to use, default is default]
* --ports[ports to use for the service.  Either a range 1-10 or a single port]

Example
---
```podlauncher --prefix=nginx --num=3 --image=nginx --namespace=default --ports="80,7000-7005"```