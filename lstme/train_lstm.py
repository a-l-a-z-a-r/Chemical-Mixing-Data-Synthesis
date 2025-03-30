import pandas as pd
import numpy as np
import pickle
from tensorflow.keras.models import Sequential, load_model
from tensorflow.keras.layers import LSTM, Dense, Dropout
from tensorflow.keras.initializers import HeNormal
from tensorflow.keras.optimizers import Adam
from sklearn.preprocessing import MinMaxScaler
from sklearn.model_selection import train_test_split

class FermentationLSTM:
    def __init__(self, csv_path, seq_length=12, epochs=50, batch_size=16, learning_rate=0.001):
        self.csv_path = csv_path
        self.seq_length = seq_length
        self.epochs = epochs
        self.batch_size = batch_size
        self.learning_rate = learning_rate
        self.model = None
        self.scaler = MinMaxScaler()
        self.feature_cols = ['time', 'biomass', 'lactic acid', 'lactose']
        self.target_col = 'pH'

    def load_data(self):
        df = pd.read_csv(self.csv_path)
        for col in self.feature_cols + [self.target_col]:
            if col not in df.columns:
                raise ValueError(f"Missing required column: {col}")

        feature_data = df[self.feature_cols].astype('float32').values
        target_data = df[self.target_col].astype('float32').values

        feature_data = self.scaler.fit_transform(feature_data)

        X, y = [], []
        for i in range(len(feature_data) - self.seq_length):
            X.append(feature_data[i:i + self.seq_length])
            y.append(target_data[i + self.seq_length])

        X = np.array(X)
        y = np.array(y)
        return train_test_split(X, y, test_size=0.2, random_state=42)

    def build_model(self, input_shape):
        model = Sequential()
        model.add(LSTM(128, return_sequences=True, input_shape=input_shape,
                       kernel_initializer=HeNormal()))
        model.add(Dropout(0.3))
        model.add(LSTM(64, return_sequences=False, kernel_initializer=HeNormal()))
        model.add(Dropout(0.3))
        model.add(Dense(32, activation='relu', kernel_initializer=HeNormal()))
        model.add(Dense(1))

        optimizer = Adam(learning_rate=self.learning_rate)
        model.compile(optimizer=optimizer, loss='mean_squared_error')
        self.model = model

    def train(self):
        X_train, X_test, y_train, y_test = self.load_data()
        self.build_model((self.seq_length, len(self.feature_cols)))
        self.model.fit(X_train, y_train, validation_data=(X_test, y_test),
                       epochs=self.epochs, batch_size=self.batch_size)

        return X_test, y_test

    def save_model(self, model_path='model.h5', scaler_path='scaler.pkl'):
        self.model.save(model_path)
        with open(scaler_path, 'wb') as f:
            pickle.dump(self.scaler, f)
        print(f"Model and scaler saved: {model_path}, {scaler_path}")

    def run(self):
        X_test, y_test = self.train()
        self.save_model()
        preds = self.model.predict(X_test).flatten()
        print("\nSample predictions:")
        for i in range(min(5, len(preds))):
            print(f"True pH: {y_test[i]:.3f} | Predicted: {preds[i]:.3f}")

if __name__ == "__main__":
    model = FermentationLSTM(csv_path="fermentation_X_0.137.csv")
    model.run()
