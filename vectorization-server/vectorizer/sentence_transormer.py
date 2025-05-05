from sentence_transformers import SentenceTransformer
import logging
from config import SENTENCE_MODEL

model = SentenceTransformer(SENTENCE_MODEL)
logging.basicConfig(level=logging.INFO)

def convert_text_to_vector(text: str):
    try:
        vector = model.encode(text,convert_to_tensor=True)
        return vector.cpu().numpy()
    except Exception as e:
        logging.error(f"Error generating vector for text: {e}")
        return None
    