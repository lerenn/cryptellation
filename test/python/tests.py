from gateway.info import *

def run_gateway_tests():
    client = Client(
        base_url="http://localhost:7003/v1",
    )
    run_system_info_test(client=client)

if __name__ == "__main__":
    run_gateway_tests()
    print("Everything passed")
