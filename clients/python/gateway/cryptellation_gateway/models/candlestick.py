from collections.abc import Mapping
from typing import Any, TypeVar, Union

from attrs import define as _attrs_define
from attrs import field as _attrs_field

from ..types import UNSET, Unset

T = TypeVar("T", bound="Candlestick")


@_attrs_define
class Candlestick:
    """
    Attributes:
        time (str): Open time of the candlestick
        open_ (float): Open price of the candlestick
        high (float): High price of the candlestick
        low (float): Low price of the candlestick
        close (float): Close price of the candlestick
        volume (float): Volume of the candlestick
        uncomplete (Union[Unset, bool]): Indicates if the candlestick is uncomplete
    """

    time: str
    open_: float
    high: float
    low: float
    close: float
    volume: float
    uncomplete: Union[Unset, bool] = UNSET
    additional_properties: dict[str, Any] = _attrs_field(init=False, factory=dict)

    def to_dict(self) -> dict[str, Any]:
        time = self.time

        open_ = self.open_

        high = self.high

        low = self.low

        close = self.close

        volume = self.volume

        uncomplete = self.uncomplete

        field_dict: dict[str, Any] = {}
        field_dict.update(self.additional_properties)
        field_dict.update(
            {
                "time": time,
                "open": open_,
                "high": high,
                "low": low,
                "close": close,
                "volume": volume,
            }
        )
        if uncomplete is not UNSET:
            field_dict["uncomplete"] = uncomplete

        return field_dict

    @classmethod
    def from_dict(cls: type[T], src_dict: Mapping[str, Any]) -> T:
        d = dict(src_dict)
        time = d.pop("time")

        open_ = d.pop("open")

        high = d.pop("high")

        low = d.pop("low")

        close = d.pop("close")

        volume = d.pop("volume")

        uncomplete = d.pop("uncomplete", UNSET)

        candlestick = cls(
            time=time,
            open_=open_,
            high=high,
            low=low,
            close=close,
            volume=volume,
            uncomplete=uncomplete,
        )

        candlestick.additional_properties = d
        return candlestick

    @property
    def additional_keys(self) -> list[str]:
        return list(self.additional_properties.keys())

    def __getitem__(self, key: str) -> Any:
        return self.additional_properties[key]

    def __setitem__(self, key: str, value: Any) -> None:
        self.additional_properties[key] = value

    def __delitem__(self, key: str) -> None:
        del self.additional_properties[key]

    def __contains__(self, key: str) -> bool:
        return key in self.additional_properties
