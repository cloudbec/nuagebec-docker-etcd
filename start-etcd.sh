#!/bin/sh


if [ -z $ETCD_DISCOVERY ]
then
    cat > /data/.env_etcd << EOF
$ETCD_DISCOVERY
EOF
fi

ETH0_IPV4=`ip -o -4 addr show | awk -F '[ /]+' '/global/ {print $4}' | head -n1`
export `cat docker_env_file` # Supposing we have only one variable. And, we want discovery URL
echo "ETCD discovery URL "$ETCD_DISCOVERY

trap 'pkill etcd' TERM

if  (etcdctl -C $ETH0_IPV4:2379 member list | grep $ETH0_IPV4)
then
        echo "this node is a cluster member"
    export ETCD_DISCOVERY=""

    while read env_var; do
        export $env_var
    done < /data/.env_etcd

    echo $(etcd -name `hostname` -listen-client-urls  http://$ETH0_IPV4:2379,http://$ETH0_IPV4:4001 -advertise-client-urls http://$ETH0_IPv4:2379,http://$ETH0_IPv4:4001  -listen-peer-urls http://$ETH0_IPV4:2380 -initial-advertise-peer-urls http://$ETH0_IPV4:2380 2>&1)
       
else
    /usr/local/bin/etcd  -name `hostname` -initial-advertise-peer-urls http://$ETH0_IPV4:2380  -listen-peer-urls http://$ETH0_IPV4:2380  -discovery $ETCD_DISCOVERY   -listen-client-urls http://$ETH0_IPV4:2379,http://$ETH0_IPV4:4001 -advertise-client-urls http://$ETH0_IPV4:2379,http://$ETH0_IPV4:4001  2>&1 1>/tmp/out.log &

    sleep 5

    if  (grep "etcd: discovery cluster full, falling back to proxy"  `find  /var/log/supervisor/ | grep etcd-stdout`)
    then
        echo "Change proxy for member"        

        sleep 3

        etcdctl -C $ETH0_IPV4:2379 member add `hostname` http://$ETH0_IPV4:2380 | grep -v "Added member named" | grep -v '^$' > /data/.env_etcd
        #  supervisorctl stop etcd


k
        export ETCD_DISCOVERY=""

        supervisorctl stop etcd

        while read env_var; do
            export $env_var
        done < /data/.env_etcd
        
        
        echo $(etcd -name `hostname` -listen-client-urls  http://$ETH0_IPV4:2379,http://$ETH0_IPV4:4001 -advertise-client-urls http://$ETH0_IPv4:2379,http://$ETH0_IPv4:4001  -listen-peer-urls http://$ETH0_IPV4:2380 -initial-advertise-peer-urls http://$ETH0_IPV4:2380 2>&1)

        supervisorctl start etcd


    else

        echo "initiate cluster"
        wait $(pgrep etcd) #hold on... since cluster run in background

    fi
fi

echo "shuting down etcd"

