#!/bin/bash


POSTGRES_CONTAINER="infrastructure-postgres-1"
KAFKA_CONNECT_CONTAINER="infrastructure-kconnect-1"
DEBEZIUM_USER="debezium_user"
DEBEZIUM_PASSWORD="debezium_user_1234"
DATABASE="productdb"
TABLE="products"
PUBLICATION_NAME="productdb_pub"
CONNECTOR_NAME="productdb-connector"


echo "Configuring PostgreSQL for logical replication..."
docker exec -it $POSTGRES_CONTAINER bash -c "echo \"wal_level = logical\" >> /var/lib/postgresql/data/postgresql.conf"
docker exec -it $POSTGRES_CONTAINER bash -c "echo \"max_replication_slots = 4\" >> /var/lib/postgresql/data/postgresql.conf"
docker exec -it $POSTGRES_CONTAINER bash -c "echo \"max_wal_senders = 4\" >> /var/lib/postgresql/data/postgresql.conf"
docker restart $POSTGRES_CONTAINER
echo "PostgreSQL configuration complete."

#  psql -U postgres -d productdb
echo "Setting up publication and replication user in PostgreSQL..."
docker exec -i $POSTGRES_CONTAINER psql -U postgres -d $DATABASE <<EOF
CREATE ROLE $DEBEZIUM_USER WITH LOGIN PASSWORD '$DEBEZIUM_PASSWORD';
ALTER ROLE $DEBEZIUM_USER SET wal_level = 'logical';
ALTER USER $DEBEZIUM_USER WITH REPLICATION;
GRANT REPLICATION ON DATABASE $DATABASE TO $DEBEZIUM_USER;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO $DEBEZIUM_USER;
CREATE PUBLICATION $PUBLICATION_NAME FOR TABLE $TABLE;
EOF
echo "Publication and user setup complete."


echo "Creating Debezium connector configuration..."
cat <<EOF > productdb_connector.json
{
  "name": "$CONNECTOR_NAME",
  "config": {
    "connector.class": "io.debezium.connector.postgresql.PostgresConnector",
    "tasks.max": "1",
    "database.hostname": "$POSTGRES_CONTAINER",
    "database.port": "5432",
    "database.user": "$DEBEZIUM_USER",
    "database.password": "$DEBEZIUM_PASSWORD",
    "database.dbname": "$DATABASE",
    "database.server.name": "${DATABASE}_server",
    "table.include.list": "public.$TABLE",
    "plugin.name": "pgoutput",
    "publication.name": "$PUBLICATION_NAME",
    "snapshot.mode": "initial"
  }
}
EOF
echo "Connector configuration created."


echo "Deploying Debezium connector..."
curl -i -X POST -H "Accept:application/json" -H "Content-Type: application/json" \
    http://localhost:8083/connectors/ -d @productdb_connector.json

echo "Debezium connector deployment complete."


echo "Checking connector status..."
sleep 3
curl -s http://localhost:8083/connectors/$CONNECTOR_NAME/status | jq
echo "Setup complete. Connector status displayed above."
# curl -s http://localhost:8083/connectors/productdb-connector/status