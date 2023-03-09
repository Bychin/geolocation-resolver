CREATE KEYSPACE geoloc
    WITH replication = {
        'class': 'NetworkTopologyStrategy',
        'DC1': 1,
        'DC2': 1,
        'DC3': 1
    }
AND durable_writes = true;
