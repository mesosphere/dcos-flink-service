#!/bin/bash
set -e
set -x

export FLINK_JOBMANAGER_WEB_PORT="$PORT0"
export FLINK_JOBMANAGER_RPC_PORT="$PORT1"
export FLINK_BLOB_SERVER_PORT="$PORT2"
export FLINK_MESOS_ARTIFACT_SERVER_PORT="$PORT3"
export LIBPROCESS_PORT="$PORT4"

export FLINK_UI_WEB_PROXY_BASE="/service/${DCOS_SERVICE_NAME}"

# start service
exec runsvdir -P /etc/service
