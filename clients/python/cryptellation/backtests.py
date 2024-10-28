
from .utils import cryptellation
from .displayer import Displayer

class Backtest(object):
    '''Backtest class'''

    def _load_data(self):
        return cryptellation(['backtests', 'get', self.id])
    
    def _load_candlesticks(self):
        return cryptellation(['candlesticks', 'read',
            '--exchange', self.data['tick_subscriptions'][0]['exchange'],
            '--pair', self.data['tick_subscriptions'][0]['pair'],
            '--period', self.data['period_between_events'],
            '--start', self.data['start_time'],
            '--end', self.data['end_time']
            ])

    def __init__(self, id):
        self.id = id
        self.data = self._load_data()

    def analyze(self):
        '''Analyze the backtest'''
        self.candlesticks = self._load_candlesticks()
        self.orders = self.list_orders()
        Displayer(self.candlesticks, self.orders).display()

    def list_orders(self):
        '''List all orders'''
        return cryptellation(['backtests', 'orders', 'list', self.id])

def list():
    '''List all backtests'''
    return cryptellation(['backtests', 'list'])

def last():
    '''Get the last backtest'''
    return Backtest(list()[-1]['id'])