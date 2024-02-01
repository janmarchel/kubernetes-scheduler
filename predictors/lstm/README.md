
# Project Title

Time-Series Forecasting with LSTM and InfluxDB

## Description

This project utilizes Long Short-Term Memory (LSTM) neural networks to forecast time-series data. It integrates with InfluxDB to store and retrieve time-series forecasts, making it suitable for applications requiring high-performance time-series data handling.

## Installation

1. Clone this repository to your local machine.
2. Ensure you have Python 3.x installed.
3. Install required dependencies by running `pip install -r requirements.txt`.

## Usage

To run the forecasting model:

1. Update `config.py` with your InfluxDB settings and model parameters.
2. Place your time-series data file in the project directory.
3. Modify `main.py` to point to your data file and specify the column to forecast.
4. Execute `python main.py` to start the forecasting process.
