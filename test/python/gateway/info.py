from cryptellation_gateway import Client
from cryptellation_gateway.api.default import get_info

def run_system_info_test(client : Client):
    info = get_info.sync(client=client)
    assert(info.version != "")