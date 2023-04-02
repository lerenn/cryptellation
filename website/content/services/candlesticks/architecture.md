---
title: "Architecture"
date: 2023-03-27T16:33:26+02:00
---

### Container Diagram

This diagram is the overview of the candlesticks service in the Cryptellation System.

```mermaid
stateDiagram-v2
    direction LR
    
    % Class definition
    classDef outside stroke:grey

    User: <center><b>User</b>\n[Application]\nService user</center>
    User --> Candlesticks: Requests candlesticks [NATS]

    Candlesticks: <center><b>Candlesticks</b>\n[Cryptellation Service]\nProvides candlesticks to\nservices and users.</center>
    Candlesticks --> Database: Reads from\nWrites to\n[Redis/SQL]
    Candlesticks --> Exchanges: Get candlesticks\n[Exchange API]

    Database: <center><b>Database</b>\n[Redis/SQL]\nStore backtests information.</center>

    Exchanges: <center><b>Exchanges</b>\n[Internet Service]\nAPIs that exchanges provides to get\ninformation and execute orders.</center>

    % Class application
    class User outside
    class Exchanges outside
```

### Component Diagram

This diagram is the internal view of the candlesticks service:

```mermaid
stateDiagram-v2
    % Class definition
    classDef outside stroke:grey
    classDef adapter stroke:green
    classDef queryCommand stroke:red
    classDef domain stroke:yellow

    User: <center><b>User</b>\n[Application]\nService user</center>
    User --> NATSController: Get request [NATS]

    Database: <center><b>Database</b>\n[Redis/SQL]\nStore backtests information.</center>

    NATSController: <center><b>NATS Controller</b>\n[Controller]\nReceives NATS request and redirect\nthem to the correct commands/queries.</center>
    NATSController --> CachedReadCandlesticks: Request candlesticks

    CachedReadCandlesticks: <center><b>Cached Read Candlesticks</b>\n[Command]\nRead candlesticks from DB,\nor from exchange API.</center>
    CachedReadCandlesticks --> DatabaseAdapter: Read candlestics [SQL]
    CachedReadCandlesticks --> CandlesticksDomain: Process data
    CachedReadCandlesticks --> ExchangesAdapter: Get candlesticks

    CandlesticksDomain: <center><b>Candlesticks Domain</b>\n[Domain]\nBusiness logic\nrelated to candlesticks.</center>

    DatabaseAdapter: <center><b>Database Adapter</b>\n[Adapter]\nAdapter to read/write DB\nthrough existing libraries and\nfor mocking purposes</center>
    DatabaseAdapter --> Database: Reads from\nWrite to\n[Redis/SQL]

    ExchangesAdapter: <center><b>Exchanges Adapter</b>\n[Adapter]\nAdapter to read candlesticks\nfrom Exchanges APIs</center>
    ExchangesAdapter --> Exchanges: Request candlesticks through API

    Exchanges: <center><b>Exchanges</b>\n[Internet Service]\nAPIs that exchanges provides to get\ninformation and execute orders.</center>

    % Class application
    class Database outside
    class User outside
    class Exchanges outside
    class NATSController adapter
    class DatabaseAdapter adapter
    class CachedReadCandlesticks queryCommand
    class CandlesticksDomain domain
```
