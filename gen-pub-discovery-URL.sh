#/bin/sh
# If env. var. named ETCD_DISCOVERY isn't set, discover it

#ETCD_DISCOVERY=$(wget --quiet https://discovery.etcd.io/new?size=$1  -O - )

#Local etcd discuvery url using unix date
VAR=`date +%s`
ETCD_DISCOVERY=http://172.17.42.1:4002/v2/keys/discovery/$VAR
curl -X PUT $ETCD_DISCOVERY/_config/size -d value=$1  #configure cluster size


cat << EOF > docker_env_file
ETCD_DISCOVERY=$ETCD_DISCOVERY
EOF

echo "etcd url is $ETCD_DISCOVERY"
echo "Env. file generated name's ./docker_env_file"
echo ###
echo "Please run \"docker run --env-file=docker_env_file nuagebec/etcd\""
