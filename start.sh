#bin.bash
	
sudo docker compose -f ./infrastructure/docker-compose.yml up --build -d
sudo docker compose -f ./entities-server/docker-compose.yml up -d
sudo docker compose -f ./recommendation-server/docker-compose.yml up  -d
sudo docker compose -f ./vectorization-server/docker-compose.yml up  -d
sudo docker compose -f ./relation-builder/docker-compose.yml up  -d
sudo ./infrastructure/setup_debezium_connector.sh
