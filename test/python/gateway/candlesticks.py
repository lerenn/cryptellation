from cryptellation_gateway import Client
from cryptellation_gateway.api.default import get_candlesticks

def run_candlesticks_test(client : Client):
    info = get_candlesticks.sync(
        client=client,
        exchange="binance",
        symbol="BTC-USDT",
        interval="M1",
        start_time="2024-01-01T00:00:00Z",
        end_time="2024-01-01T00:59:00Z",
    )
    assert(len(info) == 60)

    # Check the first candlestick
    assert(info[0].time == "2024-01-01T00:00:00Z")
    assert(info[0].open_ == 42283.58)
    assert(info[0].high == 42298.62)
    assert(info[0].low == 42261.02)
    assert(info[0].close == 42298.61)
    assert(info[0].volume == 35.92724)