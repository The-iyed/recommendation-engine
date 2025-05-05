-- Create the replication user
CREATE ROLE debezium_user WITH LOGIN PASSWORD 'debezium_user_1234' REPLICATION;
GRANT CONNECT ON DATABASE productdb TO debezium_user;

-- Create publication
CREATE PUBLICATION productdb_pub FOR TABLE public.products;
