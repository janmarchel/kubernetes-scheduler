import numpy as np
import pandas as pd
from sklearn.preprocessing import MinMaxScaler

from predictors.lstm.custom_influxdb_client import InfluxDBClientWrapper
from utils import create_sequences  # Import create_sequences from utils.py
from model import build_lstm_model  # Import build_lstm_model from model.py
import config  # Import configuration settings

class LSTMForecaster:
    def __init__(self, file_path, column_name, sequence_length=100, epochs=5, batch_size=1):
        self.file_path = file_path
        self.column_name = column_name
        self.sequence_length = sequence_length
        self.epochs = epochs
        self.batch_size = batch_size
        self.model = None
        self.scaler = MinMaxScaler(feature_range=(0, 1))

    def load_and_preprocess_data(self):
        df = pd.read_csv(self.file_path)
        self.data = self.scaler.fit_transform(df[self.column_name].values.reshape(-1, 1))

    def run(self):
        self.load_and_preprocess_data()
        X, y = create_sequences(self.data, self.sequence_length)
        train_size = int(len(X) * 0.8)
        X_train, X_test = X[:train_size], X[train_size:]
        y_train, y_test = y[:train_size], y[train_size:]
        X_train = np.reshape(X_train, (X_train.shape[0], X_train.shape[1], 1))
        X_test = np.reshape(X_test, (X_test.shape[0], X_test.shape[1], 1))

        self.model = build_lstm_model((X_train.shape[1], 1))
        self.model.fit(X_train, y_train, epochs=self.epochs, batch_size=self.batch_size, verbose=2)

        predictions = self.model.predict(X_test)
        predictions = self.scaler.inverse_transform(predictions)

        influx_client = InfluxDBClientWrapper()
        for i in range(len(predictions)):
            influx_client.write_forecast('LSTM_forecast', 'value', float(predictions[i][0]), pd.Timestamp('now'))
        influx_client.close()
