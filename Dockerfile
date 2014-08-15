FROM nuagebec/ubuntu:14.04
MAINTAINER Michaël Faille <michael.faille@nuagebec.ca>
# Let's install go and etcd (from source).


ADD ./etcd /opt/etcd

RUN  apt-get update -q && \
      DEBIAN_FRONTEND=noninteractive apt-get install -qy build-essential  && \
       apt-get clean && \
       rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/* && \  
       wget -O - https://storage.googleapis.com/golang/go1.3.1.src.tar.gz | tar -v -C /usr/local -xzi && cd /usr/local/go/src && \
       ./make.bash --no-clean 2>&1 && \
       cd /opt/etcd && PATH=/usr/local/go/bin:$PATH ./build && mv bin/* /usr/local/bin && \
       rm -R /opt/etcd && \
       apt-get remove -y --purge autoconf build-essential

ENV PATH /usr/local/go/bin:$PATH

EXPOSE 4001 7001

ENTRYPOINT ["/opt/etcd/bin/etcd"]
