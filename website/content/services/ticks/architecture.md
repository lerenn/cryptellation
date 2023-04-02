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
    User --> Ticks: Register/unregister symbol listener [NATS]

    Ticks: <center><b>Ticks</b>\n[Cryptellation Service]\nProvides ticks in real time to\nservices and users.</center>
    Ticks --> Database: Reads from\nWrites to\n[Redis/SQL]
    Ticks --> User: Deliver new tick\n[NATS]

    Database: <center><b>Database</b>\n[Redis/SQL]\nStore ticks listener\ninformation.</center>

    Exchanges: <center><b>Exchanges</b>\n[Internet Service]\nAPIs that exchanges provides to get\ninformation.</center>
    Exchanges --> Ticks: Deliver new ticks\n[HTTP WebSocket]

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

    Database: <center><b>Database</b>\n[Redis/SQL]\nStore ticks subscriptions\n information.</center>

    NATSController: <center><b>NATS Controller</b>\n[Controller]\nReceives NATS request and redirect\nthem to the correct commands/queries.</center>
    NATSController --> RegisterSymbolListener: Register\nsymbol listener
    NATSController --> UnregisterSymbolListener: Unregister\nsymbol listener

    UnregisterSymbolListener: <center><b>Unregister Symbol Listener</b>\n[Command]\nDecrease the number of listener on a\nsymbol. Terminate the listening it it\nwas the last listener.</center>
    UnregisterSymbolListener --> DatabaseAdapter: Read/Write\ninformation
    
    RegisterSymbolListener: <center><b>Register Symbol Listener</b>\n[Command]\nIncrease the number of listener on a\nsymbol. Initiate the listening if it is the\nfirst listener.</center>
    RegisterSymbolListener --> EventBrokerAdapter: Publish new ticks
    RegisterSymbolListener --> DatabaseAdapter: Read/Write\ninformation

    DatabaseAdapter: <center><b>Database Adapter</b>\n[Adapter]\nAdapter to read/write DB\nthrough existing libraries and\nfor mocking purposes</center>
    DatabaseAdapter --> Database: Reads from\nWrite to\n[Redis/SQL]

    EventBrokerAdapter: <center><b>Event Broker Adapter</b>\n[Adapter]\nAdapter to deliver events through\nexisting libraries and for\nmocking purposes.</center>
    EventBrokerAdapter --> User: Send new ticks

    ExchangesAdapter: <center><b>Exchanges Adapter</b>\n[Adapter]\nAdapter to get new ticks\nfrom Exchanges APIs</center>
    ExchangesAdapter --> RegisterSymbolListener: Send new ticks

    Exchanges: <center><b>Exchanges</b>\n[Internet Service]\nAPIs that exchanges provides to get\ninformation and execute orders.</center>
    Exchanges --> ExchangesAdapter: Send new ticks\n[HTTP WebSocket]

    % Class application
    class Database outside
    class User outside
    class Exchanges outside
    class NATSController adapter
    class DatabaseAdapter adapter
    class EventBrokerAdapter adapter
    class ExchangesAdapter adapter
    class RegisterSymbolListener queryCommand
    class UnregisterSymbolListener queryCommand
    class CandlesticksDomain domain
```
