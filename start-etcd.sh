#!/bin/sh
/usr/local/bin/etcd  -name `hostname` -initial-advertise-peer-urls https://`ip -o -4 addr show | awk -F '[ /]+' '/global/ {print $4}' | head -n1`:2380  -listen-peer-urls https:`ip -o -4 addr show | awk -F '[ /]+' '/global/ {print $4}' | head -n1`:2380  -discovery $DISCOVERY-URL
