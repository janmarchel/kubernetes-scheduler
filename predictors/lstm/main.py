# main.py

import config
from predictors.lstm.custom_influxdb_client import InfluxDBClientWrapper
from predictors.lstm.lstm import LSTMForecaster


def main():
    # Initialize LSTMForecaster with the path to your data and the column to be forecasted
    forecaster = LSTMForecaster('path_to_your_data.csv', 'YourColumnName', sequence_length=config.SEQUENCE_LENGTH, epochs=config.EPOCHS, batch_size=config.BATCH_SIZE)

    # Run the forecasting process which includes loading and preprocessing data,
    # creating sequences, training the model, making predictions, and writing to InfluxDB
    forecaster.run()

    print("Model training and forecasting complete.")

if __name__ == '__main__':
    main()


    # Initialize InfluxDB client
    influx_client = InfluxDBClientWrapper()

    # Write predictions to InfluxDB
    for i, prediction in enumerate():
        influx_client.write_forecast('lstm_forecasts', 'prediction', prediction[0], 'time_placeholder')

    # Close the InfluxDB client
    influx_client.close()

    print("Model training and forecasting complete.")

if __name__ == '__main__':
    main()