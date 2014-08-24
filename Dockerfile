FROM nuagebec/ubuntu:14.04
MAINTAINER Michaël Faille <michael.faille@nuagebec.ca>
# Let's install go and etcd (from source).


ADD supervisor-etcd.conf /etc/supervisor/conf.d/etcd.conf

RUN  apt-get update -q && \
    DEBIAN_FRONTEND=noninteractive apt-get install -qy build-essential  --no-install-recommends && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* && \
    wget -O - https://storage.googleapis.com/golang/go1.3.1.src.tar.gz | tar -v -C /usr/local -xzi && cd /usr/local/go/src && \
    ./make.bash --no-clean 2>&1 && \
    mkdir -p /data/persistent/gopath && \
    cd /opt && git clone --depth=1 https://github.com/coreos/etcd.git && cd /opt/etcd && \
    PATH=/usr/local/go/bin:$PATH ./build && mv bin/* /usr/local/bin && \
    rm -R /opt/etcd && \
    cd /opt && git clone  --depth=1 https://github.com/coreos/etcdctl/  && cd /opt/etcdctl && \
    PATH=/usr/local/go/bin:$PATH ./build && mv bin/* /usr/local/bin && \
    rm -R /opt/etcdctl && \  
    apt-get remove -y --purge build-essential

ENV PATH /usr/local/go/bin:$PATH
ENV GOPATH /data/persistent/gopath

EXPOSE 4001 7001

CMD ["/data/run.sh"]
