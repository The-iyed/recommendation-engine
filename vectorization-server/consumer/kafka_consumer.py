import json
import logging
from kafka import KafkaConsumer
from config import KAFKA_CREATE_VECTOR_TOPIC
from vectorizer.vectorizer import convert_to_vector  
from vectorizer.sentence_transormer import convert_text_to_vector
from producer.kafka_producer import send_message

consumer = KafkaConsumer(
    bootstrap_servers='kafka:29092',
    group_id='vectorization_server',
    value_deserializer=lambda m: json.loads(m.decode('utf-8'))
)
consumer.subscribe([KAFKA_CREATE_VECTOR_TOPIC])

logging.basicConfig(level=logging.INFO)


def consume_messages():
    logging.info("Starting kafka consumer")
    for message in consumer:
        try:
            product_data = message.value
            logging.info(f"Consumed message: {json.dumps(product_data, indent=2)}")
            
            image_path = product_data.get("image_path")
            fields_letter = product_data.get("fields_letter").split(".")
            
            logging.info(f"Field Letter{fields_letter}")
            product_vector = convert_to_vector(image_url=image_path)

            if product_vector is not None:
                product_data["vector"] = product_vector
            else:
                logging.warning("Failed to generate image vector.")
                continue 
            
            for field in fields_letter:
                vector_letter = convert_text_to_vector(text=product_data[field])
                if vector_letter is not None:
                    product_data[f"{field}_vector"] = vector_letter
                else:
                    logging.warning(f"Failed to generate vector for field '{field}'.")
            
            send_message(product_data)
            logging.info("Message with vectors sent back to Kafka.")


        except Exception as e:
            logging.error(f"Error processing message: {e}")