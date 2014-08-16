nuagebec-etcd
====================

etcd image for docker build on nuagebec/ubuntu:14.04 image

A highly-available key value store for shared configuration and service discovery.
etcd is inspired by Apache ZooKeeper and doozer, with a focus on being:

* *Simple*: curl'able user facing API (HTTP+JSON)
* *Secure*: optional SSL client cert authentication
* *Fast*: benchmarked 1000s of writes/s per instance
* *Reliable*: properly distributed using Raft

etcd is written in Go and uses the [Raft][raft] consensus algorithm to manage a highly-available replicated log.

Use etcdctl for a simple command line client.
Or, feel free to just use curl, as in the example below :

https://coreos.com/docs/distributed-configuration/getting-started-with-etcd/


### Related links :

* zookeeper: http://zookeeper.apache.org/
* doozer: https://github.com/ha/doozerd
* raft: http://raftconsensus.github.io/
* etcdctl: http://github.com/coreos/etcdctl/


Volume
------

/data/persitent to persist


Ports
-----

4001
7001
Running

You can use this base box standalone doing:

```bash
docker pull nuagebec/etcd:latest
docker run -d -p 4001:4001 -p 7001:7001 nuagebec/etcd:latest
```
