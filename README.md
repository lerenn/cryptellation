# Cryptellation

Cryptellation is a **scalable cryptocurrency investment system**.

This system allows developers to create bots to manage their investments on 
different cryptographic markets, featuring **backtesting**, **forward testing** and 
**live running**.

## Supported clients

* Python (documentation incoming...)

## Services 

| Service          | Description                             |
| ---------------- | --------------------------------------- |
| **Backtests**    | Execute backtests                       |
| **Candlesticks** | Get cached informations on candlesticks |
| **Exchanges**    | Get cached informations on exchanges    |
| **Livetests**    | Execute livetests                       |
| **Ticks**        | Get ticks from exchanges                |


## Running Python example

### Requirements

* bash
* docker
* docker-compose
* make
* pip

### How to

First launch the cryptellation system:

    # Copy the credentials.example.env and modify it
    cp .credentials.example.env .credentials.env

    # Launch the system
    make docker/run

Then you can use the client to execute an example:

    # Go into python client directory 
    cd clients/python

    # Install requirements and client
    pip install -r requirements.txt
    pip install -e .

    # Launch example
    python examples/graph.py # Or any other from examples/ directory
