import influxdb_client
from influxdb_client.client.write_api import SYNCHRONOUS
from influxdb_client import InfluxDBClient, Point, WritePrecision
import config

class InfluxDBClientWrapper:
    def __init__(self):
        """Initialize the InfluxDB client with configuration from config.py."""
        self.client = InfluxDBClient(url=config.INFLUXDB_URL, token=config.INFLUXDB_TOKEN, org=config.INFLUXDB_ORG)
        self.write_api = self.client.write_api(write_options=SYNCHRONOUS)

    def write_forecast(self, measurement, field, value, time):
        """Write a forecast data point to the InfluxDB."""
        p = Point(measurement).field(field, value).time(time, WritePrecision.NS)
        self.write_api.write(bucket=config.INFLUXDB_BUCKET, record=p)

    def query_forecasts(self, query):
        """Query forecasts from the InfluxDB."""
        result = self.client.query_api().query(org=config.INFLUXDB_ORG, query=query)
        forecasts = []
        for table in result:
            for record in table.records:
                forecasts.append((record.get_time(), record.get_value()))
        return forecasts

    def close(self):
        """Close the InfluxDB client connection."""
        self.client.close()