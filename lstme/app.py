from flask import Flask, request, jsonify, render_template, send_file, redirect, url_for
import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
import io
import pickle
from tensorflow.keras.models import load_model
from sklearn.metrics import r2_score, mean_absolute_percentage_error

class PHPredictor:
    def __init__(self, model_path="model.h5", scaler_path="scaler.pkl", seq_length=12):
        self.model = load_model(model_path)
        with open(scaler_path, "rb") as f:
            self.scaler = pickle.load(f)
        self.seq_length = seq_length

    def preprocess(self, df):
        required_cols = ['time', 'biomass', 'lactic acid', 'lactose', 'pH']
        for col in required_cols:
            if col not in df.columns:
                raise ValueError(f"Missing column: {col}")

        features = df[['time', 'biomass', 'lactic acid', 'lactose']].values.astype('float32')
        scaled_features = self.scaler.transform(features)

        sequences = []
        for i in range(len(scaled_features) - self.seq_length + 1):
            sequences.append(scaled_features[i:i + self.seq_length])

        sequences = np.array(sequences)
        actual_ph = df['pH'].values[self.seq_length - 1:]
        time_points = df['time'].values[self.seq_length - 1:]
        return sequences, actual_ph, time_points

    def predict(self, sequences):
        preds = self.model.predict(sequences)
        return preds.flatten()


class PHPredictorAPI:
    def __init__(self, app: Flask, predictor: PHPredictor):
        self.app = app
        self.predictor = predictor
        self.results_table = []
        self.performance_metrics = {}

        self.setup_routes()

    def setup_routes(self):
        self.app.route('/', methods=['GET'])(self.index)
        self.app.route('/upload', methods=['POST'])(self.upload_file)
        self.app.route('/results', methods=['GET'])(self.results)
        self.app.route('/plot.png', methods=['GET'])(self.plot_png)

    def index(self):
        return render_template('index.html')

    def upload_file(self):
        if 'file' not in request.files or request.files['file'].filename == '':
            return redirect(url_for('index'))

        try:
            file = request.files['file']
            df = pd.read_csv(file)
            sequences, actual_ph, time = self.predictor.preprocess(df)
            predicted_ph = self.predictor.predict(sequences)
        except Exception as e:
            print(f"Upload error: {e}")
            return redirect(url_for('index'))

        alpha_error = [abs(a - p) for a, p in zip(actual_ph, predicted_ph)]
        self.results_table = [
            {
                "Time": round(float(t), 3),
                "Actual pH": round(float(a), 3),
                "Predicted pH": round(float(p), 3),
                "Alpha Error": round(float(err), 3)
            }
            for t, a, p, err in zip(time, actual_ph, predicted_ph, alpha_error)
        ]

        #Performance Metrics
        r2 = r2_score(actual_ph, predicted_ph)
        mape = mean_absolute_percentage_error(actual_ph, predicted_ph) * 100  # %
        self.performance_metrics = {
            "R² Score": round(r2, 4),
            "MAPE (%)": round(mape, 2),
            "Accuracy (%)": round(100 - mape, 2)
        }

        return redirect(url_for('results'))

    def results(self):
        # Ensure both data and metrics exist before rendering
        if not self.results_table or not self.performance_metrics:
            return redirect(url_for('index'))
        
        # Explicitly define fallback values
        metrics = self.performance_metrics or {
            "R² Score": "N/A",
            "MAPE (%)": "N/A",
            "Accuracy (%)": "N/A"
        }

        return render_template('results.html', table=self.results_table, metrics=metrics)



    def plot_png(self):
        time = [entry['Time'] for entry in self.results_table]
        actual = [entry['Actual pH'] for entry in self.results_table]
        predicted = [entry['Predicted pH'] for entry in self.results_table]

        fig, ax = plt.subplots(figsize=(10, 5), facecolor="#181818")
        ax.set_facecolor("#202020")

        ax.plot(time, actual, label="Actual pH", color="#00BFFF", linestyle="dashed", linewidth=2)
        ax.plot(time, predicted, label="Predicted pH", color="#FF4500", linewidth=2)
        ax.fill_between(time, actual, predicted, color="gray", alpha=0.3, label="Alpha (Error)")

        ax.set_xlabel("Time", color="white")
        ax.set_ylabel("pH Value", color="white")
        ax.set_title("Predicted vs Actual pH", color="white")
        ax.legend()
        ax.grid(color="gray")

        for spine in ax.spines.values():
            spine.set_color("white")
        plt.xticks(color="white")
        plt.yticks(color="white")

        img_io = io.BytesIO()
        plt.savefig(img_io, format='png', facecolor=fig.get_facecolor())
        img_io.seek(0)
        return send_file(img_io, mimetype='image/png')

if __name__ == '__main__':
    app = Flask(__name__)
    predictor = PHPredictor()
    api = PHPredictorAPI(app, predictor)
    app.run(debug=True)
