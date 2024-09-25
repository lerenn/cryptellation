import sys
import json
import pandas as pd
import mplfinance as mpf

class Displayer(object):
    def __init__(self, data):
        self.candlesticks = data["candlesticks"]
        self.orders = data["orders"]

    def display(self):
        for exchange in self.candlesticks:
            for pair in self.candlesticks[exchange]:
                cds = self._load_candlesticks(self.candlesticks[exchange][pair])
                buy, sell = self._load_orders(self.orders)
                cds['buy'] = buy['price']
                cds['sell'] = sell['price']
                mpf.plot(cds, type='candle', warn_too_much_data=100000, addplot=[
                    mpf.make_addplot(cds['buy'], scatter=True, marker='^'),
                    mpf.make_addplot(cds['sell'], scatter=True, marker='v')
                ])
                
    def _load_candlesticks(self, candlesticks):
        df = pd.DataFrame(candlesticks)
        df['time'] = pd.to_datetime(df['time'])
        df.set_index('time', inplace=True)
        return df
    
    def _load_orders(self, orders):
        df = pd.DataFrame(orders)
        df['time'] = pd.to_datetime(df['execution_time'])
        df.set_index('time', inplace=True)

        return df[df['side'] == 'buy'], df[df['side'] == 'sell']

if __name__ == '__main__':
    with open(sys.argv[1]) as f:
        d = json.load(f)

    displayer = Displayer(d)
    displayer.display()