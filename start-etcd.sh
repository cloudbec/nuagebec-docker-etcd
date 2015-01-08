#!/bin/sh
ETH0_IPv4=`ip -o -4 addr show | awk -F '[ /]+' '/global/ {print $4}' | head -n1`

/usr/local/bin/etcd  -name `hostname` -initial-advertise-peer-urls https://$ETH0_IPv4:2380  -listen-peer-urls https://$ETH0_IPv4:2380  -discovery $ETCD_DISCOVERY   -listen-client-urls http://$ETH0_IPv4:2379,http://$ETH0_IPv4:4001
