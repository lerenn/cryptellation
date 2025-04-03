from http import HTTPStatus
from typing import Any, Optional, Union

import httpx

from ... import errors
from ...client import AuthenticatedClient, Client
from ...models.candlestick import Candlestick
from ...types import UNSET, Response, Unset


def _get_kwargs(
    *,
    exchange: str,
    symbol: str,
    interval: str,
    start_time: Union[Unset, str] = UNSET,
    end_time: Union[Unset, str] = UNSET,
) -> dict[str, Any]:
    params: dict[str, Any] = {}

    params["exchange"] = exchange

    params["symbol"] = symbol

    params["interval"] = interval

    params["start_time"] = start_time

    params["end_time"] = end_time

    params = {k: v for k, v in params.items() if v is not UNSET and v is not None}

    _kwargs: dict[str, Any] = {
        "method": "get",
        "url": "/candlesticks",
        "params": params,
    }

    return _kwargs


def _parse_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Optional[list["Candlestick"]]:
    if response.status_code == 200:
        response_200 = []
        _response_200 = response.json()
        for response_200_item_data in _response_200:
            response_200_item = Candlestick.from_dict(response_200_item_data)

            response_200.append(response_200_item)

        return response_200
    if client.raise_on_unexpected_status:
        raise errors.UnexpectedStatus(response.status_code, response.content)
    else:
        return None


def _build_response(
    *, client: Union[AuthenticatedClient, Client], response: httpx.Response
) -> Response[list["Candlestick"]]:
    return Response(
        status_code=HTTPStatus(response.status_code),
        content=response.content,
        headers=response.headers,
        parsed=_parse_response(client=client, response=response),
    )


def sync_detailed(
    *,
    client: Union[AuthenticatedClient, Client],
    exchange: str,
    symbol: str,
    interval: str,
    start_time: Union[Unset, str] = UNSET,
    end_time: Union[Unset, str] = UNSET,
) -> Response[list["Candlestick"]]:
    """Get candlestick data

    Args:
        exchange (str):
        symbol (str):
        interval (str):
        start_time (Union[Unset, str]):
        end_time (Union[Unset, str]):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[list['Candlestick']]
    """

    kwargs = _get_kwargs(
        exchange=exchange,
        symbol=symbol,
        interval=interval,
        start_time=start_time,
        end_time=end_time,
    )

    response = client.get_httpx_client().request(
        **kwargs,
    )

    return _build_response(client=client, response=response)


def sync(
    *,
    client: Union[AuthenticatedClient, Client],
    exchange: str,
    symbol: str,
    interval: str,
    start_time: Union[Unset, str] = UNSET,
    end_time: Union[Unset, str] = UNSET,
) -> Optional[list["Candlestick"]]:
    """Get candlestick data

    Args:
        exchange (str):
        symbol (str):
        interval (str):
        start_time (Union[Unset, str]):
        end_time (Union[Unset, str]):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        list['Candlestick']
    """

    return sync_detailed(
        client=client,
        exchange=exchange,
        symbol=symbol,
        interval=interval,
        start_time=start_time,
        end_time=end_time,
    ).parsed


async def asyncio_detailed(
    *,
    client: Union[AuthenticatedClient, Client],
    exchange: str,
    symbol: str,
    interval: str,
    start_time: Union[Unset, str] = UNSET,
    end_time: Union[Unset, str] = UNSET,
) -> Response[list["Candlestick"]]:
    """Get candlestick data

    Args:
        exchange (str):
        symbol (str):
        interval (str):
        start_time (Union[Unset, str]):
        end_time (Union[Unset, str]):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        Response[list['Candlestick']]
    """

    kwargs = _get_kwargs(
        exchange=exchange,
        symbol=symbol,
        interval=interval,
        start_time=start_time,
        end_time=end_time,
    )

    response = await client.get_async_httpx_client().request(**kwargs)

    return _build_response(client=client, response=response)


async def asyncio(
    *,
    client: Union[AuthenticatedClient, Client],
    exchange: str,
    symbol: str,
    interval: str,
    start_time: Union[Unset, str] = UNSET,
    end_time: Union[Unset, str] = UNSET,
) -> Optional[list["Candlestick"]]:
    """Get candlestick data

    Args:
        exchange (str):
        symbol (str):
        interval (str):
        start_time (Union[Unset, str]):
        end_time (Union[Unset, str]):

    Raises:
        errors.UnexpectedStatus: If the server returns an undocumented status code and Client.raise_on_unexpected_status is True.
        httpx.TimeoutException: If the request takes longer than Client.timeout.

    Returns:
        list['Candlestick']
    """

    return (
        await asyncio_detailed(
            client=client,
            exchange=exchange,
            symbol=symbol,
            interval=interval,
            start_time=start_time,
            end_time=end_time,
        )
    ).parsed
