from dotenv import load_dotenv
import os

load_dotenv()

KAFKA_BOOTSTRAP_SERVERS = os.getenv("KAFKA_BOOTSTRAP_SERVERS")
KAFKA_CREATE_VECTOR_TOPIC = os.getenv("KAFKA_CREATE_VECTOR_TOPIC")
KAFKA_VECTOR_CREATED_TOPIC = os.getenv("KAFKA_VECTOR_CREATED_TOPIC")
VECTOR_SIZE = int(os.getenv("VECTOR_SIZE", 512))  
HOST = os.getenv("HOST", "0.0.0.0")
PORT = int(os.getenv("PORT", 5003))  
SENTENCE_MODEL = os.getenv("SENTENCE_MODEL")
