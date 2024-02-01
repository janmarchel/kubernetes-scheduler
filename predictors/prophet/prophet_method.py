import pandas as pd
from prophet import Prophet


class ProphetForecaster:
    def __init__(self, data_path, target_column, additional_regressors=None):
        """
        Initializes the Prophet forecaster.

        :param data_path: Path to the CSV file containing the time series data.
        :param target_column: Name of the column in the CSV file to forecast.
        :param additional_regressors: Optional list of names of additional regressor columns.
        """
        self.data_path = data_path
        self.target_column = target_column
        self.additional_regressors = additional_regressors
        self.model = Prophet()

        if additional_regressors:
            for regressor in additional_regressors:
                self.model.add_regressor(regressor)

    def run(self):
        """
        Loads the data, prepares it for Prophet, and fits the model.
        """
        df = pd.read_csv(self.data_path)

        # Ensure the dataframe has the correct column names for Prophet ('ds' for dates and 'y' for values)
        if 'ds' not in df.columns:
            # Assuming the date or datetime column is named 'Date'
            df.rename(columns={'Date': 'ds', self.target_column: 'y'}, inplace=True)
        else:
            df.rename(columns={self.target_column: 'y'}, inplace=True)

        # Fit the model
        self.model.fit(df)

    def predict(self, periods=30, freq='D'):
        """
        Generates future dates and makes predictions.

        :param periods: Number of periods to forecast into the future.
        :param freq: Frequency of the forecast ('D' for days, 'W' for weeks, etc.)
        :return: A dataframe with the forecast components.
        """
        future = self.model.make_future_dataframe(periods=periods, freq=freq)
        forecast = self.model.predict(future)
        return forecast[['ds', 'yhat', 'yhat_lower', 'yhat_upper']]