# Includes libmesos.so
FROM mesosphere/mesos:1.0.1-2.0.93.ubuntu1404

RUN apt-get update && \
    apt-get install -y \
            runit \
            curl

# The base image contains java 7, but it has no environment variables set for it.
ENV JAVA_HOME /usr/lib/jvm/java-7-openjdk-amd64/jre

ENV MESOS_NATIVE_JAVA_LIBRARY /usr/lib/libmesos.so

ENV HADOOP_CONF_DIR /etc/hadoop

RUN mkdir /etc/hadoop

ADD runit/service /var/lib/runit/service

ADD runit/init.sh /sbin/init.sh

RUN chmod +x /sbin/init.sh

RUN ln -s /var/lib/runit/service/flink /etc/service/flink

RUN chmod -R ugo+rw /etc/service

RUN chmod -R ugo+rw /var/lib/

RUN chmod -R ugo+rw /var/run/

RUN chmod -R ugo+rw /var/log/

WORKDIR /

# Copy custom build to image.
COPY flink/flink-dist/target/flink-1.2-SNAPSHOT-bin/ .

# Copy base Flink configuration to image.
COPY conf/ flink-1.2-SNAPSHOT/conf/

ENV FLINK_HOME /flink-1.2-SNAPSHOT

ENV FLINK_CONF_DIR /flink-1.2-SNAPSHOT/conf

WORKDIR flink-1.2-SNAPSHOT
