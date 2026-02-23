# Diagram Alur SNMP-ZTE

File ini berisi kumpulan diagram Mermaid yang menjelaskan alur kerja sistem SNMP-ZTE. Anda dapat melihat diagram ini menggunakan preview Markdown yang mendukung Mermaid, atau menyalin kodenya ke [Mermaid Live Editor](https://mermaid.live/).

## 1. Arsitektur Global

```mermaid
graph TD
    Client((User/Client)) -->|HTTP Request| Handlers[API Layer: Handlers]
    Handlers -->|Method Calls| Services[Service Layer: Services]
    Services -->|Get/Set| Cache[(Cache Layer: Redis)]
    Services -->|Interface Implementation| Drivers[Hardware Layer: Drivers]
    Drivers -->|SNMP Get/Walk| OLT[[ZTE OLT Device]]

    subgraph "Core Logic"
        Handlers
        Services
        Drivers
    end
```

## 2. Alur Startup Aplikasi

```mermaid
sequenceDiagram
    participant Main as cmd/api/main.go
    participant Config as internal/config
    participant Redis as internal/cache
    participant Svc as internal/service
    participant Router as setupRouter()

    Main->>Config: Load() (Baca olts.json)
    Main->>Redis: NewClient() (Koneksi Redis)
    Main->>Svc: NewONUService & NewOLTService
    Svc->>Svc: Inisialisasi Driver untuk tiap OLT
    Main->>Router: Definisikan Endpoint & Middleware
    Main->>Main: ListenAndServe() (Server Aktif)
```

## 3. Alur Request: Query Stateless (/api/v1/query)

```mermaid
sequenceDiagram
    participant User
    participant QHandler as QueryHandler
    participant Driver as ZTE Driver
    participant SNMP as SNMP Client

    User->>QHandler: POST /api/v1/query (JSON Body)
    QHandler->>QHandler: getDriver() (C320)
    QHandler->>Driver: Connect()
    Driver->>SNMP: Establish UDP Connection
    QHandler->>Driver: GetONUList()
    Driver->>SNMP: SNMP WALK (OID .1.3.6.1.4.1.3902...)
    SNMP-->>Driver: SNMP Packets
    Driver-->>QHandler: List ONUInfo
    QHandler-->>User: JSON Response (Data + Duration)
```

## 4. Alur Request: Query Berbasis ID (/api/v1/olts/{id}/...)

```mermaid
sequenceDiagram
    participant User
    participant OHandler as ONUHandler
    participant Service as ONUService
    participant Cache as Redis Cache
    participant Driver as ZTE Driver

    User->>OHandler: GET /api/v1/olts/OLT-01/board/1/pon/1
    OHandler->>Service: GetONUList(oltID, board, pon)
    Service->>Cache: Get(cacheKey)
    alt Cache Ada (Hit)
        Cache-->>Service: Data ONU (JSON)
    else Cache Kosong (Miss)
        Service->>Driver: GetONUList()
        Driver->>Driver: SNMP Walk to OLT
        Driver-->>Service: List ONUInfo
        Service->>Cache: Set(cacheKey, Data)
    end
    Service-->>OHandler: Result
    OHandler-->>User: JSON Response
```
