#!/usr/bin/env bash

the_env=$1
if [[ "$the_env" == "" ]]; then
    the_env="dev"
fi

echo "env: ${the_env}";

part=1;
repl=1;
zookaddr="localhost:2181";

topics=(

"topic1.${the_env}"

"topic2.${the_env}"

)

for i in ${topics[*]}; do
  echo "create topic: ${i}"
/data/kafka/bin/kafka-topics.sh \
--create \
--zookeeper ${zookaddr} \
--replication-factor ${repl} \
--partitions ${part} \
--topic ${i}
done