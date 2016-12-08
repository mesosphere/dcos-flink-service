# dcos-flink-service

The DC/OS service repository for Apache Flink.

This repository consists of two parts:
- The docker container for the Apache Flink AppMaster/JobManager.
- The service package.

This repository contains submodules, so use following commands to clone the
repository:

```
git clone --recursive https://github.com/mesosphere/dcos-flink-service
```

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
