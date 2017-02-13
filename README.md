# dcos-flink-service

**This is experimental!**

The DC/OS service repository for Apache Flink.

This repository consists of two parts:
- The docker container for the Apache Flink AppMaster/JobManager.
- The service package.

This repository contains submodules, so use following commands to clone the
repository:

```
git clone --recursive https://github.com/mesosphere/dcos-flink-service
```

# Testing Apache Flink

Prequisite:
- DC/OS 1.8 cluster
- DC/OS CLI installed.

1. Install local package repository.
  1. `dcos marathon app add https://raw.githubusercontent.com/mesosphere/dcos-flink-service/master/service/marathon-universe.json` and wait until deployment has finished.
  2. `dcos package repo add --index=0 dev-universe http://universe.marathon.mesos:8085/repo`
2. Install Flink service
  1. `dcos package install flink`
3. Access UI
  1. You can access the UI via services via <cluster name>/service/flink/. Unfortunately you cannot upload a job jar via UI in this case.
  2. Use marathon-lb `dcos package install marathon lb` and update the marathon app definition with the following labels:

~~~
  "labels":{
    "HAPROXY_GROUP":"external",
    "HAPROXY_0_VHOST":"<public host name, e.g., ELB>"
  }
~~~

3. Submit Job

  1. UI

  2. Docker Container

      `dcos node ssh --leader --master-proxy`

      `./bin/flink run -m <jobmangerhost>:<jobmangerjobmanager.rpc.port> ./examples/batch/WordCount.jar  --input file:///etc/resolv.conf --output file:///etc/wordcount_out`

    Note that the input/output is local to the container, so most likely you want to configure HDFS or S3.


## AppMaster Docker Container

The main DC/OS service launches on this container. It can be easily built from
the root of this repository by executing

```
make -C container/appmaster
```

> **NOTE**
>
> The AppMaster image requires currently an own master build of Flink.
> Compiling Flink takes about 10-15 minutes. If you are working with this
> repository and have already built Flink once, you can invoke
> `make -C container/appmaster build-container` which only rebuilds the docker
> image (this works as long as you don't make changes to Flink or upgrade it,
> then you need a full rebuild).

## Service Package

The service package contains all needed files to make up an own DC/OS universe
package.

To try them out, you need a running DC/OS cluster and a correctly setup `dcos`
binary pointing to your cluster. Also you need to create your own universe
service that will serve the Flink package on your cluster. Refer to the
instructions
[here](https://github.com/mesosphere/universe/blob/version-3.x/README.md) to
set up your own universe server.

## Todos

* CLI support (in progress)
* HA setup (requires HDFS...)
* Update to release version once it is available
* Add more examples

