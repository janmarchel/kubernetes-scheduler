import pandas as pd
from matplotlib import pyplot as plt
from prophet.plot import add_changepoints_to_plot


def prepare_holidays(holidays_file):
    """
    Loads holiday data from a file.

    :param holidays_file: Path to the CSV file containing holidays.
    :return: A DataFrame with 'holiday' and 'ds' columns.
    """
    holidays_df = pd.read_csv(holidays_file)
    return holidays_df

def plot_with_changepoints(forecast, model, figsize=(10, 6)):
    """
    Plots the forecast along with change points.

    :param forecast: The forecast DataFrame returned by the predict method.
    :param model: The Prophet model instance.
    :param figsize: Tuple specifying figure size.
    """
    fig = model.plot(forecast, figsize=figsize)
    add_changepoints_to_plot(fig.gca(), model, forecast)
    plt.show()

def add_holidays(df, holidays_df):
    df = pd.merge(df, holidays_df, on='ds', how='left')
    df['holiday'] = df['holiday'].fillna(0)  # Assuming 'holiday' column in holidays_df
    return df