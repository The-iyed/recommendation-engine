from kafka import KafkaProducer
import logging
import json
from config import KAFKA_VECTOR_CREATED_TOPIC
import numpy as np


logging.basicConfig(level=logging.INFO)

producer = KafkaProducer(
    bootstrap_servers='kafka:29092',
    value_serializer=lambda v: json.dumps(v).encode('utf-8')
)

def send_message(product):
    try:

        if isinstance(product.get("vector"), np.ndarray):
            product["vector"] = product["vector"].tolist()
            
        for key, value in product.items():
            if key.endswith("_vector") and isinstance(value, np.ndarray):
                product[key] = value.tolist()
                
        product_data = {
            "id": product["product_id"],
            "product_id": product["product_id"],
            "name": product["name"],
            "description": product["description"],
            "price": product["price"],
            "image_path": product["image_path"],
            "vector": product["vector"],
            "fields_letter": product["fields_letter"],
        }
        for key, value in product.items():
            if key.endswith("_vector"):
                product_data[key] = value

        logging.info(f"Sending product message: {json.dumps(product_data, indent=2)}")
        
        producer.send(KAFKA_VECTOR_CREATED_TOPIC, value=product_data)
        producer.flush()  
        
        logging.info(f"Message sent to topic {KAFKA_VECTOR_CREATED_TOPIC}")
    
    except Exception as e:
        logging.error(f"Error sending product data to Kafka: {e}")
    
def close_producer():
    producer.close()