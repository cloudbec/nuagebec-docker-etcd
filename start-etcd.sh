#!/bin/sh
ETH0_IPv4=`ip -o -4 addr show | awk -F '[ /]+' '/global/ {print $4}' | head -n1`
trap 'pkill etcd' EXIT

/usr/local/bin/etcd  -name `hostname` -initial-advertise-peer-urls http://$ETH0_IPv4:2380  -listen-peer-urls http://$ETH0_IPv4:2380  -discovery $ETCD_DISCOVERY   -listen-client-urls http://$ETH0_IPv4:2379,http://$ETH0_IPv4:4001 -advertise-client-urls http://$ETH0_IPv4:2379,http://$ETH0_IPv4:4001

