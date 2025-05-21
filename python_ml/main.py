from fastapi import FastAPI, File, UploadFile
from landmark_detection import predict_landmark
app = FastAPI()

@app.post("/predict")
async def predict(file: UploadFile = File(...)):
    image_data = await file.read()
    result = predict_landmark(image_data)
    return {"label": result["label"],"address":  result["address"], "lat": result["lat"], "lon": result["lon"]}
