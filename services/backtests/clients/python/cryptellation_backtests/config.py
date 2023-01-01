from typing import List


class Config:
    def __init__(
        self,
        pubsub_url: str = "127.0.0.1:4222",
        service_url: str = "127.0.0.1:9004",
    ):
        self.pubsub_url = pubsub_url
        self.service_url = service_url
