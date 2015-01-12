#!/bin/sh

ETH0_IPv4=`ip -o -4 addr show | awk -F '[ /]+' '/global/ {print $4}' | head -n1`

trap 'pkill etcd' TERM

if etcdctl -C ETH0_IPv4:2379 member list | grep ETH0_IPv4
then

    /usr/local/bin/etcd  -name `hostname` -initial-advertise-peer-urls http://$ETH0_IPv4:2380  -listen-peer-urls http://$ETH0_IPv4:2380  -discovery $ETCD_DISCOVERY   -listen-client-urls http://0.0.0.0:2379,http://0.0.0.0:4001 -advertise-client-urls http://$ETH0_IPv4:2379,http://$ETH0_IPv4:4001

else    
    if [ -x /data/`hostname`.etcd ] 
    then

        export ETCD_DISCOVERY=""

        while read env_var; do
            export $env_var
        done < /data/.env_etcd

        etcd -name `hostname` -listen-client-urls  http://0.0.0.0:2379,http://0.0.0.0:4001 -advertise-client-urls http://$ETH0_IPv4:2379,http://$ETH0_IPv4:4001  -listen-peer-urls http://$ETH0_IPV4:2380 -initial-advertise-peer-urls http://$ETH0_IPV4:2380
        
    else

        /usr/local/bin/etcd  -name `hostname` -initial-advertise-peer-urls http://$ETH0_IPv4:2380  -listen-peer-urls http://$ETH0_IPv4:2380  -discovery $ETCD_DISCOVERY   -listen-client-urls http://0.0.0.0:2379,http://0.0.0.0:4001 -advertise-client-urls http://$ETH0_IPv4:2379,http://$ETH0_IPv4:4001 &

        sleep 10
        
        etcdctl -C ETH0_IPv4:2379 member add `hostname` http://ETH0_IPv4:2380 > /data/.env_etcd
        supervisorctl stop etcd
        supervisorctl start etcd

    fi
fi


