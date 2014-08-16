nuagebec-etcd
====================


etcd image for docker build on nuagebec/ubuntu:14.04 image

Volume
====================

/data/persitent to persist


Ports
====================

4001
7001
Running

You can use this base box standalone doing:

docker pull nuagebec/etcd:latest
    docker run -d -p 4001:4001 -p 7001:7001 nuagebec/etcd:latest
