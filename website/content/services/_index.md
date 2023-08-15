+++
archetype = "chapter"
title = "Services"
weight = 3
+++

On this section, you will find all information to understand what every service
is and how they are interacting with each other.

{{% children description="true" %}}

## System Diagrams

According to the [C4model methodology](https://c4model.com/), you'll find every
diagram to represent how every service is placed in the whole Cryptellation ecosystem. 

### System Context Diagram

This diagram is the first step to understanding the big picture of what's
happening under the hood of the Cryptellation system.

```mermaid
stateDiagram-v2
    direction LR

    % Class definition
    classDef user stroke:yellow
    classDef outside stroke:grey

    User: <center><b>User</b>\n[Application]\nService user</center>
    User --> Cryptellation: Execute real and\nsimulated orders

    Cryptellation: <center><b>Cryptellation</b>\n[Software System]\nExecute backtests, forward tests\nand live runs of strategies</center>
    Cryptellation --> User: Get information
    Cryptellation --> Exchanges: Execute real orders
    Cryptellation --> Cryptellation: Run simulations

    Exchanges: <center><b>Exchanges</b>\n[Internet Service]\nAPIs that exchanges provides to get\ninformation and execute orders</center>
    Exchanges --> Cryptellation: Get information

    % Class application
    class User user
    class Exchanges outside
```

### System Container Diagram

This diagram is the complete overview of every services in the Cryptellation System.

```mermaid
stateDiagram-v2
   classDef user stroke:yellow
   classDef outside stroke:grey

    User: <center><b>User</b>\n[Application]\nService user.</center>
    User --> CryptellationBacktests: Send orders
    User --> CryptellationForwardTests: Send orders
    User --> CryptellationExchanges: Get exchanges information
    User --> CryptellationCandlesticks: Get candlesticks
    User --> CryptellationIndicators: Get indicators

    CryptellationIndicators: <center><b>Indicators</b>\n[Cryptellation Service]\nProcess and cache common indicators.</center>
    CryptellationIndicators --> User: Deliver calculated/cached indicators
    CryptellationIndicators --> CryptellationCandlesticks: Get candlesticks

    CryptellationBacktests: <center><b>Backtests</b>\n[Cryptellation Service]\nProvide an entire environment\nfor users to execute backtests.</center>
    CryptellationBacktests --> User: Deliver simulated ticks
    CryptellationBacktests --> CryptellationCandlesticks: Get candlesticks

    CryptellationForwardTests: <center><b>Forward Tests</b>\n[Cryptellation Service]\nProvide an entire environment\nfor users to execute forward tests.</center>

    CryptellationCandlesticks: <center><b>Candlesticks</b>\n[Cryptellation Service]\nProvide cached candlesticks\nhistory to services and users.</center>
    CryptellationCandlesticks --> Exchanges: Get candlesticks

    CryptellationExchanges: <center><b>Exchanges</b>\n[Cryptellation Service]\nProvide cached exchanges\ninformation to services and users.</center>
    CryptellationExchanges --> Exchanges: Get information

    CryptellationTicks: <center><b>Ticks</b>\n[Cryptellation Service]\nProvide ticks in real time\nto services and users.</center>
    CryptellationTicks --> User: Deliver new ticks
    CryptellationTicks --> CryptellationForwardTests: Deliver new ticks

    Exchanges: <center><b>Exchanges</b>\n[Internet Service]\nAPIs that exchanges provides to get\ninformation and execute orders.</center>
    Exchanges --> CryptellationTicks: Deliver new ticks

    class User user
    class Exchanges outside
```

## Other diagrams

For more information on services and more detail diagrams, you should go in each
and every service subdirectory.

