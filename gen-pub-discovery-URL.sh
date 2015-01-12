#/bin/sh
# If env. var. named ETCD_DISCOVERY isn't set, discover it

ETCD_DISCOVERY=$(wget --quiet https://discovery.etcd.io/new?size=$1  -O - )

cat << EOF > docker_env_file
ETCD_DISCOVERY=$ETCD_DISCOVERY
EOF

echo "etcd url is $ETCD_DISCOVERY"
echo "Env. file generated name's ./docker_env_file"
echo ###
echo "Please run \"docker run --env-file=docker_env_file nuagebec/etcd\""
