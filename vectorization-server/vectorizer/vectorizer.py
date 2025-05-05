from PIL import Image
import torch
from torchvision import models, transforms
import requests
from io import BytesIO
import logging

model = models.resnet50(weights=models.ResNet50_Weights.IMAGENET1K_V1)
model.eval()
 
transform = transforms.Compose([
    transforms.Resize((224, 224)),
    transforms.ToTensor(),
    transforms.Normalize(mean=[0.485, 0.456, 0.406], std=[0.229, 0.224, 0.225]),
])
logging.basicConfig(level=logging.INFO)

def convert_to_vector(image_url: str):
    try:
        response = requests.get(image_url)
        response.raise_for_status()  
        img = Image.open(BytesIO(response.content)).convert("RGB")
        img_tensor = transform(img).unsqueeze(0)
        with torch.no_grad():
            vector = model(img_tensor).flatten()
        return vector.numpy() 

    except requests.exceptions.RequestException as e:
        logging.error(f"Failed to fetch image from URL: {e}")
        return None
    except Exception as e:
        logging.error(f"Error generating vector: {e}")
        return None
