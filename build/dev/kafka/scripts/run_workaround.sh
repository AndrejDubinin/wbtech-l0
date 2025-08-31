# This script is required to run kafka cluster (without zookeeper)
#!/bin/sh

echo "Start run_workaround.sh"

# Docker workaround: Remove check for KAFKA_ZOOKEEPER_CONNECT parameter
sed -i '/KAFKA_ZOOKEEPER_CONNECT/d' /etc/confluent/docker/configure

# Docker workaround: Ignore cub zk-ready
sed -i 's/cub zk-ready/echo ignore zk-ready/' /etc/confluent/docker/ensure

# KRaft required step: Format the storage directory with a new cluster ID
sed -i "/kafka-storage format/d" /etc/confluent/docker/ensure

cluster_id=$(kafka-storage random-uuid)

echo -e "\nkafka-storage format --ignore-formatted -t $cluster_id -c /etc/kafka/kafka.properties" >> /etc/confluent/docker/ensure
