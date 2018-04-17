Redis is generically described as a in-memory database.  This document is
not intended to describe how to use redis, but rather how this operator
manages redis.

# Types of Redis Solutions

## Traditional Redis

Traditional Redis, as opposed to Redis Cluster, is the more simplistic redis
situation
* each instance has a complete copy of the database
* one instance handles all writes
* redis sentinel is used for high availability solutions

This is the redis that one is typically given in the default [AWS Elasticache](https://aws.amazon.com/elasticache/redis/)
or [Azure Redis](https://docs.microsoft.com/en-us/azure/redis-cache/)

Its a great solution for typical workloads and is generally less restricted
in use compared to Redis Cluster.

### Redis Sentinel

In order for Traditional Redis to function in a HA mode, Redis Sentinel
is typically used.  Using Redis Sentinel we can properly handle failover

#### Sentinel Cluster

Redis Sentinel works as a cluster itself.  Using at least three nodes it
can achieve consensus if a redis node is down.  If enough sentinel nodes
agree that a node is down, a script can be executed to handle the failover.

#### Redis Pools

Redis Sentinel monitors pools of redis nodes to determine if the current
master of a redis pool is down and then elect a new master.  Because of the
nature of this problem, one sentinel cluster can monitor multiple redis pools.

#### Script to handle failover

When a redis master has failed and a new one is to be elected, a script
can be executed to handle this operation.  In a kubernetes environement,
this script can trigger the adjustment of labels to redis pods.

Consider a world where a redis service is "all pods tagged `master=true`"
The update script could ensure that only the currently-elected master has
the `master=true` tag applied.  The result is that a consumer of the service
never has to know which redis node is the master, just whatever is reached
by connecting to the redis service.

## Redis Cluster

Redis Cluster is the idea that a single redis solution may be too big for
a single machine.  Redis Cluster offers the ability to scale redis up to
1000 nodes while maintaining high performance.  There are some key concepts
to understand, however you can also read about them from the [official documentation](https://docs.microsoft.com/en-us/azure/redis-cache/)

### Hash Slots

Every key that is stored in redis will be hashed using the key name.  
There are ways to indicate which part of the key should be used for hashing
purposes, however by default the entire key name is used.  This hash will
result in a value between 0 and 16384.  The hash _is_ 16 bit, but only 14
bits are actually used to determine the slot id.

All keys with with the same hash slot will be stored on the same redis instance.

### Node Groups

A node group is a grouping of nodes within redis.  It does not have to be
more than one node, but in typical situations there are at least two nodes
per node group.  A node group will be assigned an exclusive set of of hash
slots.  This means that only one node pool will ever have the data for any
given hash slot.

### Node Group Roles

Within a node group, there will only be one master.  This master will handle
all client read and write requests for the hash slots it is given control
over.  Replicas will only replicate the master and will normally not allow
for client read/write requests as they may not be up to date, however a
replica can be configured to issue read requests.

### Restrictions

Because of the nature of redis cluster, multi-key operations can only happen
on the same hash slot.  The same goes for Lua scripts.

### Manual Failover

Redis cluster can be told to failover a master node.  This is done by having
a replica of the node group execute the `CLUSTER FAILOVER` command.  Doing
so will cause redis to gracefully elect and new master and handle the promotion.

This is the preferred action when doing redis instance maintenance.

### Master Election

The cluster, like many other distributed systems, emits heartbeats to
member nodes.  If a node group replica has discovered that its master has
gone offline, it can ask the cluster for promotion.  Remaining masters within
the cluster will decide if a promotion should happen.

### Networking

#### Intra-Redis Cluster

Within redis cluster, traffic (like heartbeats) are handled on a separate
port from client traffic.  This separate port is fixed - it is always the
result of adding 10000 to the client port.  For example, by default redis
listens on 6379.  This means that redis-cluster also communicates on
port 16379.  This also means that redis client traffic can never be higher
than 55535.

#### Client-Server Traffic

It is expected that each redis cluster is accessible to clients through their
advertised port.  This means that if there are 1000 redis nodes in a cluster
then the redis client may have to talk to any one of those hosts directly.

This _direct communication_ requirement can be an important issue when
dealing with NATs.  However this also shows how well fit kubernetes is for
redis-cluster.  With every instance having its own IP address, there are no
issues within a kubernetes cluster.

# Shared Redis Considerations

## Memory usage

Generically speaking, you can tell redis how much memory it is allowed to
use and it will obey you.  There is some overhead for the actual redis process
but generally speaking the max amount of memory being used by the process will be
close to the directive.

Redis _does_ use more memory when handling replication concerns however.

Whenever redis is performing a backup or is trying to get a replica online
it will fork itself.  The memory limit you give redis is without considering
the forking.  This means that, from a kubernetes point of view, if you want
1Gi usable memory within a redis system, you will need to provision 2Gi of memory.
Perhaps a little bit additional to handle any overhead would be prudent.

## Disk usage

As far as capacity is concerned, Redis Cluster has two types of backup solutions.  
A traditional backup file and an append only file.  In both cases, the
disk usage needed will be greater than the max allocated memory to redis.  
This is because a new backup will be created before the existing one is
removed.  This means that if you're going to have a 1Gi redis instance,
you should provision at least 2Gi of disk space for backups.

As far as I/O utilization, backups will happen frequently and as such one
should imagine a rather consistent stream of data being written to disk.  
To make matters worse, if a backup cannot be completed, redis will immediately
restart its efforts to do a backup, thus potentially causing for even higher
disk I/O than in a normal, healthy situation.

Kubernetes does not allow for the disk I/O resource limitations, however
external data providers (AWS EBS, etc.) _do_ - so be sure to be aware of
the disk I/O profile of your redis application before provisioning too slow
of disks.