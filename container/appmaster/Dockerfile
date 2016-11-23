FROM mesosphere/mesos:1.0.1-2.0.93.ubuntu1404

# The base image contains java 7, but it has no environment variables set for it.
ENV JAVA_HOME /usr/lib/jvm/java-7-openjdk-amd64/jre

# Copy custom build to image.
COPY flink/flink-dist/target/flink-1.2-SNAPSHOT-bin/ .

COPY conf/ flink-1.2-SNAPSHOT/conf/
WORKDIR flink-1.2-SNAPSHOT

# Set up files that shall appear in the $MESOS_SANDBOX later. All these files will be symlinked there.
RUN mkdir mesos_sandbox \
    && ln -s ../lib/flink-dist_2.10-1.2-SNAPSHOT.jar mesos_sandbox/flink.jar \
    && ln -s ../lib/flink-python_2.10-1.2-SNAPSHOT.jar mesos_sandbox/flink-python_2.10-1.2-SNAPSHOT.jar \
    && ln -s ../lib/log4j-1.2.17.jar mesos_sandbox/log4j-1.2.17.jar \
    && ln -s ../lib/slf4j-log4j12-1.7.7.jar mesos_sandbox/slf4j-log4j12-1.7.7.jar \
    && ln -s ../conf/flink-conf.yaml mesos_sandbox/flink-conf.yaml \
    && ln -s ../conf/log4j.properties mesos_sandbox/log4j.properties

ENV _CLIENT_SHIP_FILES flink-python_2.10-1.2-SNAPSHOT.jar,log4j-1.2.17.jar,slf4j-log4j12-1.7.7.jar,log4j.properties
ENV _FLINK_CLASSPATH *

ENV _CLIENT_TM_MEMORY 1024
ENV _CLIENT_TM_COUNT 1
ENV _SLOTS 2
ENV _CLIENT_USERNAME root
ENV _CLIENT_SESSION_ID default

CMD ln -s $(pwd)/mesos_sandbox/* $MESOS_SANDBOX/ \
    && $JAVA_HOME/bin/java -cp "lib/*" -Dlog.file=jobmaster.log -Dlog4j.configuration=file:log4j.properties org.apache.flink.mesos.runtime.clusterframework.MesosApplicationMasterRunner --configDir .
