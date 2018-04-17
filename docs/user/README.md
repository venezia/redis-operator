# User Documentation

## Deploying the Redis Operator

There are two ways of deploying the redis operator:
* Manually through provided examples
* Through the helm chart provided

The helm chart should be sufficient for most people and should be the
easiest to use.

### Operator Resource Requirements

This is still a work in progress, however these are our best estimates:

| Resource | Minimum | Suggested |
| --- | --- | --- |
| CPU | TBD | TBD |
| Memory | TBD | TBD |
| Network | TBD | TBD |
| Disk Usage | TBD | TBD |

### Privileges Required

Because the Redis Operator is going to create and maintain redis installations
it will need privileges to do so.  From an RBAC perspective, it will need
to be able to do the following:

```yaml
kind: ClusterRole
apiVersion: rbac.authorization.k8s.io/v1rules:
- apiGroups: ["redis.database.samsung.com"]
  resources: ["*"]
  verbs: ["*"]
- apiGroups: [""]
  resources: ["pods", "services", "secrets", "jobs"]
  verbs: ["*"]
- apiGroups: ["extensions", "apps"]
  resources: ["deployments"]
  verbs: ["*"]
```

Improved privilege requirements will be updated here later.

### Deploying through Helm Chart

_*Please always review the helm default values to ensure there are no changes*_

To Install:
* Add Chart Repository
* Create values file if neccessary
* execute
```helm install repo/chart-name --name redis-operator --values values.yaml```

### Verifying Installation

A couple quick ways to verify the installation of the operator is to look
for the presence of the API and the CRD.

#### Verifying the API

```shell
$ kubectl api-versions |grep redis.database
redis.database.samsung.com/v1alpha1
```

#### Verifying the CRD

```shell
$ kubectl get crd redii.redis.database.samsung.com
NAME                               AGE
redii.redis.database.samsung.com   sometime
```

## Creating a Redis Instance

In order to create a redis instance, please create a redis api object.

An example might look like

```yaml
apiVersion: redis.database.samsung.com/v1alpha1
kind: Redis
metadata:
  namespace: default
  name: my-redis
spec:
  redis:
    replicas: 3
    requests:
      memory: 100Mi
  setinel:
    replicas: 3
```

As the operator is developed further, further examples will be located in
[examples/user folder](../../examples/user) of the project