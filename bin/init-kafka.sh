#!/bin/bash

cd $DATAPATH_BIN

source ./init.sh

# Get Kafka client
go get -u github.com/jdamick/kafka

# Build client and tools
cd $GOPATH/src/github.com/jdamick/kafka
make kafka
make tools


# =====
# Kafka
# =====
cd $KAFKA_HOME

# Start Zookeeper
zookeeper-server-start.sh config/zookeeper.properties

# Start Kafka
kafka-server-start.sh config/server.properties

# Create a topic
kafka-topics.sh --create --zookeeper localhost:2181 --replication-factor 1 --partitions 1 --topic messages

# Check for topic
# bin/kafka-topics.sh --list --zookeeper localhost:2181

# Test Sending messages using {"res":"http://domain.com/test.csv"}
# bin/kafka-console-producer.sh --broker-list localhost:9092 --topic messages

# Start consumer and check for messages
# bin/kafka-console-consumer.sh --zookeeper localhost:2181 --topic messages --from-beginning
