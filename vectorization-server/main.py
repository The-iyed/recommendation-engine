from fastapi import FastAPI
import threading
import logging
from consumer.kafka_consumer import consume_messages

app = FastAPI()

@app.on_event("startup")
async def startup_event():
    consumer_thread = threading.Thread(target=consume_messages)
    consumer_thread.start()
    logging.info("Kafka consumer started.")

@app.get("/")
async def read_root():
    return {"message": "Vectorization server is running."}
