#!/bin/bash
set -e
set -x

export FLINK_JOBMANAGER_WEB_PORT="$PORT0"
export FLINK_JOBMANAGER_RPC_PORT="$PORT1"
export FLINK_BLOB_SERVER_PORT="$PORT2"
export FLINK_MESOS_ARTIFACT_SERVER_PORT="$PORT3"
export LIBPROCESS_PORT="$PORT4"

export FLINK_UI_WEB_PROXY_BASE="/service/${DCOS_SERVICE_NAME}"
export FLINK_LOG_DIR=$MESOS_SANDBOX

# validate base64 encoded keystore and truststroe
if [[ "${FLINK_SSL_ENABLED}" == true ]]; then
	KEYDIR=`mktemp -d`
	trap "rm -rf $KEYDIR" EXIT

	echo "${FLINK_SSL_KEYSTOREBASE64}" | base64 -d > "$KEYDIR/flink.keystore"
	ALIAS=$(keytool -list -keystore "$KEYDIR/flink.keystore" -storepass "${FLINK_SSL_KEYSTOREPASSWORD}" | grep PrivateKeyEntry | cut -d, -f1 | head -n1)
	if [[ -z "${ALIAS}" ]]; then
		echo "Cannot find private key in keystore"
		exit 1
	fi

	echo "${FLINK_SSL_TRUSTSTOREBASE64}" | base64 -d > "$KEYDIR/flink.truststore"
	ALIAS=$(keytool -list -keystore "$KEYDIR/flink.truststore" -storepass "${FLINK_SSL_TRUSTSTOREPASSWORD}" | grep trustedCertEntry | cut -d, -f1 | head -n1)
	if [[ -z "${ALIAS}" ]]; then
		echo "Cannot find trusted cert entry in keystore"
		exit 1
	fi

	rm -rf "$KEYDIR"
fi

# Move hadoop config files, as specified by hdfs.config-url, into place.
if [[ -f $MESOS_SANDBOX/hdfs-site.xml && -f $MESOS_SANDBOX/core-site.xml ]]; then
		echo "Copying HDFS configuration files to ${HADOOP_CONF_DIR}"
    mkdir -p "${HADOOP_CONF_DIR}"
    cp $MESOS_SANDBOX/hdfs-site.xml "${HADOOP_CONF_DIR}"
    cp $MESOS_SANDBOX/core-site.xml "${HADOOP_CONF_DIR}"
fi

# Move kerberos config file, as specified by security.kerberos.krb5conf, into place.
if [[ "${FLINK_SECURITY_KRB5_CONF_BASE64}" != "" ]]; then
    echo "${FLINK_SECURITY_KRB5_CONF_BASE64}" | base64 -d > /etc/krb5.conf
fi

# start service
exec runsvdir -P /etc/service
