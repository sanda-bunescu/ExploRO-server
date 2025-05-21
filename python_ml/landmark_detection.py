import io
import time
from geopy.geocoders import Nominatim
import tensorflow as tf
import pandas as pd
import tensorflow_hub as hub
import numpy as np
from PIL import Image

# Constants
img_shape = (321, 321)

# Load model and labels once
model_url = 'https://www.kaggle.com/models/google/landmarks/TensorFlow1/classifier-europe-v1/1'
classifier = hub.KerasLayer(model_url, output_key='predictions:logits', input_shape=img_shape + (3,))

df = pd.read_csv('landmarks_classifier_europe_V1_label_map.csv')
labels = dict(zip(df.id, df.name))

geolocator = Nominatim(user_agent="landmark-geolocator")


def predict_landmark(image_bytes):
    # Load and preprocess image
    image = Image.open(io.BytesIO(image_bytes)).resize(img_shape).convert('RGB')
    image = np.array(image) / 255.0
    image = image.astype(np.float32)[np.newaxis, ...]

    # Run inference
    logits = classifier(image).numpy()
    top_idx = int(np.argmax(logits))

    # Get landmark name
    landmark_name = labels.get(top_idx, "Unknown")

    # Try geolocation
    location = geolocator.geocode(landmark_name)
    if location:
        return {
            "label": landmark_name,
            "address": location.address,
            "lat": location.latitude,
            "lon": location.longitude
        }
    else:
        return {
            "label": landmark_name,
            "address": None,
            "lat": None,
            "lon": None
        }
