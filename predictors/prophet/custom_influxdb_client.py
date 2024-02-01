from influxdb_client import InfluxDBClient
from influxdb_client.client.write_api import SYNCHRONOUS
import config

class InfluxDBClientWrapper:
    def __init__(self):
        """Initialize the InfluxDB client with configuration from config.py."""
        self.client = InfluxDBClient(url=config.INFLUXDB_URL, token=config.INFLUXDB_TOKEN, org=config.INFLUXDB_ORG)
        self.write_api = self.client.write_api(write_options=SYNCHRONOUS)


    def write_forecast(self, measurement, forecast_df):
        for _, row in forecast_df.iterrows():
            json_body = [
                {
                    "measurement": measurement,
                    "time": row['ds'].to_pydatetime(),
                    "fields": {
                        "yhat": row['yhat'],
                        "yhat_lower": row['yhat_lower'],
                        "yhat_upper": row['yhat_upper']
                    }
                }
            ]
            self.client.write_points(json_body)

    def close(self):
        self.client.close()
