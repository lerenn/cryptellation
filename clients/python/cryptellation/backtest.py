import grpc
import threading
import queue
from typing import Dict, List
import iso8601
import nats
import asyncio
import json
from datetime import datetime

from .account import Account
from .event import Event
from .order import Order
from .period import Period
from .grapher import Grapher
from .candlesticks import Candlesticks

from cryptellation.config import Config

import cryptellation._genproto.backtests_pb2 as backtests
import cryptellation._genproto.backtests_pb2_grpc as backtests_grpc


class BacktestEvents(threading.Thread):
    def __init__(self, id, pubsub_url):
        threading.Thread.__init__(self)
        self._id = id
        self._pubsub_url = pubsub_url
        self._events_queue = queue.Queue(maxsize=0)
        self.start()

    async def _receive(self, pubsub_url):
        nc = await nats.connect(pubsub_url)
        sub = await nc.subscribe("Backtests.%d" % (self._id,))

        async for msg in sub.messages:
            self._events_queue.put(msg)

    def run(self):
        asyncio.run(self._receive(self._pubsub_url))

    def get(self) -> Event:
        e = self._events_queue.get()
        backtest_event = backtests.BacktestEvent()
        backtest_event.ParseFromString(e.data)
        return Event(
            iso8601.parse_date(backtest_event.time),
            backtest_event.type,
            json.loads(backtest_event.content),
        )


class Backtest(object):
    def __init__(
        self,
        start_time: datetime,
        end_time: datetime = datetime.now(),
        accounts: Dict[str, Account] = {
            "binance": Account({"USDC": 100000}),
        },
    ):
        self._config = Config()
        self._candlesticks = Candlesticks()
        self._channel = grpc.insecure_channel(self._config.backtests_url)
        self._stub = backtests_grpc.BacktestsServiceStub(self._channel)
        self._id = self._create(start_time, end_time, accounts)
        self._start_time = start_time
        self._actual_time = self._start_time
        self._end_time = end_time

    def _create(self, start_time, end_time, accounts) -> int:
        if start_time.tzinfo is None or start_time.tzinfo.utcoffset(start_time) is None:
            raise Exception("no timezone specified on start")

        return self._stub.CreateBacktest(
            backtests.CreateBacktestRequest(
                start_time=start_time.isoformat(),
                end_time=end_time.isoformat(),
                accounts=self._account_to_grpc(accounts),
            )
        ).id

    def _account_to_grpc(self, accounts: Dict[str, Account]):
        req_accounts = {}
        for exch, account in accounts.items():
            assets = {}
            for asset, quantity in account.assets.items():
                assets[asset] = quantity
            req_accounts[exch] = backtests.Account(assets=assets)
        return req_accounts

    def subscribe(self, exchange_name, pair_symbol):
        self._stub.SubscribeToBacktestEvents(
            backtests.SubscribeToBacktestEventsRequest(
                id=self._id,
                exchange_name=exchange_name,
                pair_symbol=pair_symbol,
            )
        )

    def advance(self):
        self._stub.AdvanceBacktest(
            backtests.AdvanceBacktestRequest(
                id=self._id,
            )
        )

    def listen(self) -> BacktestEvents:
        return BacktestEvents(self._id, self._config.pubsub_url)

    def order(self, type: str, exchange: str, pair: str, side: str, quantity: float):
        req = backtests.CreateBacktestOrderRequest(
            backtest_id=self._id,
            order = backtests.Order(
                type=type,
                exchange_name=exchange,
                pair_symbol=pair,
                side=side,
                quantity=quantity,
            )
        )
        self._stub.CreateBacktestOrder(req)

    def accounts(self) -> Dict[str, Account]:
        req = backtests.BacktestAccountsRequest(
            backtest_id=self._id,
        )
        resp = self._stub.BacktestAccounts(req)
        return self._grpc_to_accounts(resp)

    def _grpc_to_accounts(self, resp: backtests.BacktestAccountsResponse) -> Dict[str, Account]:
        accounts = {}
        for exch, account in resp.accounts.items():
            assets = {}
            for asset, quantity in account.assets.items():
                assets[asset] = quantity
            accounts[exch] = Account(assets)
        return accounts

    def orders(self) -> List[Order]:
        req = backtests.BacktestOrdersRequest(backtest_id=self._id)
        resp = self._stub.BacktestOrders(req)
        return self._grpc_orders(resp)

    def _grpc_orders(self, resp: backtests.BacktestOrdersResponse) -> List[Order]:
        orders = []
        for o in resp.orders:
            orders.append(
                Order(
                    iso8601.parse_date(o.execution_time),
                    o.type,
                    o.exchange_name,
                    o.pair_symbol,
                    o.side,
                    o.quantity,
                    o.price,
                )
            )
        return orders

    def on_event(self, event: Event):
        pass

    def on_end(self):
        pass

    def display(self, exchange: str, pair: str, period: Period):
        p = Grapher()

        start = self._start_time
        end = self._end_time
        cs = self._candlesticks.get(exchange, pair, period, start, end)
        p.candlesticks(cs)

        p.orders(self.orders())

        p.show()

    def actual_time(self) -> datetime:
        return self._actual_time

    def run(self):
        events = self.listen()
        finished = False
        while finished is False:
            self.advance()

            while True:
                event = events.get()

                if event.type == "status":
                    finished = event.content["finished"]
                    self._actual_time = event.time
                    break

                self.on_event(event)

        return self.on_end()

    def candlesticks(
        self,
        exchange: str,
        pair: str,
        period: Period,
        relative_start: int,
        relative_end: int = 0,
        limit: int = 0,
    ):
        start = self._actual_time - relative_start * period.duration()
        end = self._actual_time - relative_end * period.duration()
        return self._candlesticks.get(exchange, pair, period, start, end, limit)
