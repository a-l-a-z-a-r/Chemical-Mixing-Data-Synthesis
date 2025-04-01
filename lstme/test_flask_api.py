import io
import numpy as np
import pandas as pd
import pytest
from predictor_api import PHPredictor, PHPredictorAPI
from flask import Flask

@pytest.fixture
def client():
    app = Flask(__name__)
    predictor = PHPredictor(model_path="model.h5", scaler_path="scaler.pkl")
    PHPredictorAPI(app, predictor)
    return app.test_client()

def create_dummy_csv():
    df = pd.DataFrame({
        "time": np.arange(20),
        "biomass": np.random.rand(20),
        "lactic acid": np.random.rand(20),
        "lactose": np.random.rand(20),
        "pH": 3 + np.random.rand(20)
    })
    stream = io.StringIO()
    df.to_csv(stream, index=False)
    stream.seek(0)
    return stream

def test_index_page(client):
    response = client.get("/")
    assert response.status_code == 200

def test_file_upload_success(client):
    data = {
        "file": (create_dummy_csv(), "test.csv")
    }
    response = client.post("/upload", data=data, content_type="multipart/form-data", follow_redirects=True)
    assert response.status_code == 200
    assert b"Actual pH" in response.data

def test_plot_png(client):
    # First upload to populate data
    data = {"file": (create_dummy_csv(), "test.csv")}
    client.post("/upload", data=data, content_type="multipart/form-data")
    
    # Then get plot
    response = client.get("/plot.png")
    assert response.status_code == 200
    assert response.mimetype == "image/png"
