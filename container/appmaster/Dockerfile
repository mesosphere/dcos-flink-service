# Includes libmesos.so
FROM mesosphere/mesos:1.0.1-2.0.93.ubuntu1404

# The base image contains java 7, but it has no environment variables set for it.
ENV JAVA_HOME /usr/lib/jvm/java-7-openjdk-amd64/jre

WORKDIR /

# Copy custom build to image.
COPY flink/flink-dist/target/flink-1.2-SNAPSHOT-bin/ .
# Copy base Flink configuration to image.
COPY conf/ flink-1.2-SNAPSHOT/conf/

WORKDIR flink-1.2-SNAPSHOT

ENV FLINK_HOME /flink-1.2-SNAPSHOT
ENV FLINK_CONF_DIR /flink-1.2-SNAPSHOT/conf

CMD cd "$MESOS_SANDBOX" && "$FLINK_HOME/bin/mesos-appmaster.sh"
