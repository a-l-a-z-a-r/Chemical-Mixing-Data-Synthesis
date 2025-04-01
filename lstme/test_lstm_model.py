import os
import pytest
import numpy as np
import pandas as pd
from fermentation_lstm import FermentationLSTM

def dummy_csv(tmp_path):
    path = tmp_path / "dummy.csv"
    df = pd.DataFrame({
        "time": np.arange(20),
        "biomass": np.random.rand(20),
        "lactic acid": np.random.rand(20),
        "lactose": np.random.rand(20),
        "pH": 3 + np.random.rand(20)
    })
    df.to_csv(path, index=False)
    return str(path)

def test_load_data(dummy_csv):
    lstm = FermentationLSTM(csv_path=dummy_csv, seq_length=5)
    X_train, X_test, y_train, y_test = lstm.load_data()
    assert X_train.shape[1:] == (5, 4)
    assert len(y_train) > 0

def test_model_training(dummy_csv):
    lstm = FermentationLSTM(csv_path=dummy_csv, seq_length=5, epochs=1)
    X_test, y_test = lstm.train()
    assert lstm.model is not None
    assert X_test.shape[0] == y_test.shape[0]

def test_model_save(dummy_csv, tmp_path):
    model_path = tmp_path / "test_model.h5"
    scaler_path = tmp_path / "test_scaler.pkl"

    lstm = FermentationLSTM(csv_path=dummy_csv, seq_length=5, epochs=1)
    lstm.train()
    lstm.save_model(model_path, scaler_path)

    assert os.path.exists(model_path)
    assert os.path.exists(scaler_path)
