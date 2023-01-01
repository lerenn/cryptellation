from typing import List


class Config:
    def __init__(
        self,
        url: str = "127.0.0.1:9005",
    ):
        self.url = url
