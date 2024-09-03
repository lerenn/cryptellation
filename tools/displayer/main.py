import sys
import json
import pandas as pd
import mplfinance as mpf

class Displayer(object):
    def __init__(self, data):
        self.candlesticks = data["candlesticks"]

    def display(self):
        for exchange in self.candlesticks:
            for pair in self.candlesticks[exchange]:
                self._display_candlesticks(self.candlesticks[exchange][pair])
                
    def _display_candlesticks(self, candlesticks):
        df = pd.DataFrame(candlesticks)
        df['time'] = pd.to_datetime(df['time'])
        df.set_index('time', inplace=True)
        mpf.plot(df, type='candle', warn_too_much_data=100000)

if __name__ == '__main__':
    with open(sys.argv[1]) as f:
        d = json.load(f)

    displayer = Displayer(d)
    displayer.display()