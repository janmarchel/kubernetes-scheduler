from predictors.prophet import config
from predictors.prophet.custom_influxdb_client import InfluxDBClientWrapper
from predictors.prophet.prophet_method import ProphetForecaster


def main():
    # Initialize and run the Prophet forecaster
    forecaster = ProphetForecaster(config.DATA_PATH, config.TARGET_COLUMN)
    forecaster.run()
    forecast = forecaster.predict(periods=60)  # Adjust periods as needed

    # Initialize InfluxDB client and write forecast results
    influx_client = InfluxDBClientWrapper(**config.DATABASE)
    influx_client.write_forecast('forecast_measurement', forecast)
    influx_client.close()

    print("Forecasting and data writing complete.")

if __name__ == '__main__':
    main()