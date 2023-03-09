#!/bin/bash

# deploy_dev launches and initialises scylladb in docker.
# Note: it must be run from the root project dir.

read -r -d '' create_ks << EOM
CREATE KEYSPACE geoloc
    WITH replication = {
        'class': 'SimpleStrategy',
        'replication_factor': 1
    }
AND durable_writes = true;
EOM

create_entries=$(cat databases/scylla/table_entry.cql)

function start() {
    docker-compose -f scripts/dev/docker-compose.yml up -d

    if ! docker ps | grep -q "scylla-node"; then
        echo "scylla has not started"
        exit 1
    fi

    docker exec scylla-node1 cqlsh -e "DESC geoloc" &> /dev/null
    ret=$?

    if [[ $ret != 0 ]]; then
        echo "Initialising scylla... (may take some time)"
        sleep 120 # cause scylla takes some time before starting listening
        docker exec scylla-node2 cqlsh -e "$create_ks"       || exit 1
        docker exec scylla-node2 cqlsh -e "$create_entries"  || exit 1
        echo "Successfully initialised scylla"
    fi
}

function stop() {
    docker-compose -f scripts/dev/docker-compose.yml down
}

function restart() {
    stop
    start
}

function status() {
    if [[ $(docker ps | egrep -c "scylla-node") != 3 ]]; then
        echo "something is wrong"
        exit 1
    fi

    echo "everything seems to be ok"
}


case $1 in
    start) start;;
    stop) stop;;
    restart) restart;;
    status) status;;
    *) echo "Usage: deploy_dev.sh start|stop|restart|status";;
esac
