# 🚀 Arquitectura de Eventos en Tiempo Real - CPI Hub API

## 📋 Resumen Ejecutivo

Esta aplicación implementa un sistema de eventos en tiempo real usando WebSockets con una arquitectura limpia que separa responsabilidades entre capas. El sistema permite comunicación en tiempo real entre usuarios dentro de espacios específicos.

## 🏗️ Arquitectura General

```mermaid
graph TB
    subgraph "Frontend/Cliente"
        FE[Cliente Web/Mobile]
    end
    
    subgraph "API Layer"
        ROUTER[Router Gin]
        HANDLER[EventsHandler]
    end
    
    subgraph "Business Logic"
        USECASE[EventsUsecase]
        HUB_MGR[HubManager]
        CLIENT_MGR[ClientManager]
    end
    
    subgraph "Infrastructure"
        WS[WebSocket Wrapper]
        REPO[EventsRepository]
        DB[(PostgreSQL)]
    end
    
    subgraph "Domain"
        DOMAIN[Domain Models]
        DTO[DTOs]
    end
    
    FE -->|WebSocket Connection| ROUTER
    FE -->|HTTP POST Broadcast| ROUTER
    ROUTER --> HANDLER
    HANDLER --> USECASE
    USECASE --> HUB_MGR
    USECASE --> CLIENT_MGR
    USECASE --> REPO
    CLIENT_MGR --> WS
    HUB_MGR --> WS
    REPO --> DB
    USECASE --> DOMAIN
    HANDLER --> DTO
```

## 🔄 Flujo de Conexión WebSocket

```mermaid
sequenceDiagram
    participant C as Cliente
    participant R as Router
    participant H as EventsHandler
    participant U as EventsUsecase
    participant HM as HubManager
    participant CM as ClientManager
    participant WS as WebSocket
    
    C->>R: GET /ws/spaces/{space_id}?user_id={user_id}
    R->>H: Connect()
    H->>U: HandleConnection(params, w, r)
    U->>WS: Upgrade HTTP to WebSocket
    WS-->>U: WebSocket Connection
    U->>U: CreateClient(userID, spaceID, conn)
    U->>HM: RegisterClient(client)
    HM->>HM: Add to Clients map
    HM->>HM: Send Join message to space
    U->>CM: NewClientManager(client)
    U->>CM: Start ReadPump() goroutine
    U->>CM: Start WritePump() goroutine
    CM-->>C: Connection established
```

## 💬 Flujo de Mensajes de Chat

```mermaid
sequenceDiagram
    participant C1 as Cliente 1
    participant C2 as Cliente 2
    participant CM as ClientManager
    participant HM as HubManager
    participant R as EventsRepository
    participant DB as Database
    
    C1->>CM: Send chat message
    CM->>CM: Validate message
    CM->>CM: Create ChatMessage
    CM->>HM: Broadcast to space
    HM->>HM: Send to all clients in space
    HM->>C1: Echo message
    HM->>C2: Forward message
    HM->>R: Save message to DB
    R->>DB: INSERT chat_message
    DB-->>R: Success
    R-->>HM: Message saved
```

## 🔌 Flujo de Desconexión

```mermaid
sequenceDiagram
    participant C as Cliente
    participant CM as ClientManager
    participant HM as HubManager
    participant C2 as Otros Clientes
    
    Note over C: Cliente se desconecta<br/>(cierra navegador, pérdida red, etc.)
    CM->>CM: ReadPump() detects disconnection
    CM->>CM: defer function executes
    CM->>HM: Send client to Unregister channel
    CM->>CM: Close WebSocket connection
    HM->>HM: Remove client from Clients map
    HM->>HM: Close client.Send channel
    HM->>HM: Create Leave message
    HM->>C2: Broadcast Leave message to space
    C2-->>C2: Update UI (user left)
```

## 📁 Estructura de Archivos y Responsabilidades

### 🎯 **Domain Layer** (`internal/core/domain/`)
- **`events.go`**: Modelos de dominio (Client, Hub, EventMessage, ChatMessage, etc.)
- **`repositories.go`**: Interfaces de repositorios

### 🔧 **Use Cases** (`internal/core/usecase/events/`)
- **`events_usecase.go`**: Lógica de negocio principal
  - `HandleConnection()`: Maneja conexiones WebSocket
  - `Broadcast()`: Procesa mensajes de chat
  - `CreateClient()`: Crea clientes
  - `RegisterClient()`: Registra clientes en el hub

- **`hub_manager.go`**: Gestión del hub central
  - `Run()`: Loop principal del hub
  - `BroadcastChatMessage()`: Difunde mensajes de chat
  - `broadcastToSpace()`: Envía mensajes a espacios específicos

- **`client_manager.go`**: Gestión individual de clientes
  - `ReadPump()`: Lee mensajes del cliente
  - `WritePump()`: Escribe mensajes al cliente
  - `handleMessage()`: Procesa mensajes recibidos

### 🌐 **Infrastructure Layer**

#### **Handlers** (`internal/infrastructure/entrypoint/handlers/events/`)
- **`events.go`**: Endpoints HTTP
  - `Connect()`: Establece conexión WebSocket
  - `Broadcast()`: Envía mensajes via HTTP

#### **Adapters** (`internal/infrastructure/adapters/`)
- **`websocket/websocket_wrapper.go`**: Wrapper para Gorilla WebSocket
- **`repositories/postgres/events/events_repository.go`**: Persistencia de mensajes

## 🔄 Estados del Sistema

### 1. **Inicialización**
```mermaid
graph LR
    A[App Start] --> B[Create HubManager]
    B --> C[Start Hub.Run() goroutine]
    C --> D[Hub Ready]
```

### 2. **Conexión de Cliente**
```mermaid
graph LR
    A[WebSocket Request] --> B[Upgrade Connection]
    B --> C[Create Client]
    C --> D[Register in Hub]
    D --> E[Start ClientManager]
    E --> F[Client Active]
```

### 3. **Manejo de Mensajes**
```mermaid
graph LR
    A[Message Received] --> B[Validate Message]
    B --> C[Process Message]
    C --> D[Broadcast to Space]
    D --> E[Save to Database]
    E --> F[Notify Clients]
```

### 4. **Desconexión**
```mermaid
graph LR
    A[Connection Lost] --> B[Detect Disconnection]
    B --> C[Unregister Client]
    C --> D[Clean Resources]
    D --> E[Notify Other Clients]
```

## 🚦 Tipos de Mensajes

| Tipo | Propósito | Cuándo se Envía |
|------|-----------|-----------------|
| `join` | Notificar entrada de usuario | Al conectarse |
| `leave` | Notificar salida de usuario | Al desconectarse |
| `chat` | Mensaje de chat | Al enviar mensaje |
| `ping` | Mantener conexión viva | Periódicamente |
| `pong` | Respuesta a ping | En respuesta a ping |
| `error` | Mensaje de error | En caso de error |

## 🔧 Configuración y Timeouts

```go
const (
    writeWait      = 10 * time.Second  // Timeout para escribir
    pongWait       = 60 * time.Second  // Timeout para pong
    pingPeriod     = 54 * time.Second  // Intervalo de ping
    maxMessageSize = 512               // Tamaño máximo de mensaje
)
```

## 🎯 Endpoints Disponibles

### WebSocket
- **`GET /v1/ws/spaces/{space_id}?user_id={user_id}`**
  - Establece conexión WebSocket
  - Parámetros: `space_id` (path), `user_id` (query)

### HTTP REST
- **`POST /v1/ws/spaces/{space_id}/broadcast`**
  - Envía mensaje de chat via HTTP
  - Body: `{"message": "texto", "user_id": "id", "username": "name"}`

## 🔒 Consideraciones de Seguridad

1. **Validación de Parámetros**: Todos los parámetros son validados
2. **Límites de Mensaje**: Tamaño máximo de 512 bytes
3. **Timeouts**: Conexiones se cierran automáticamente si no responden
4. **Limpieza de Recursos**: Recursos se liberan automáticamente

## 📊 Métricas y Logging

El sistema incluye logging detallado para:
- Conexiones y desconexiones
- Errores de WebSocket
- Mensajes enviados/recibidos
- Problemas de difusión

## 🚀 Escalabilidad

La arquitectura actual está diseñada para:
- **Múltiples espacios**: Cada espacio maneja sus propios clientes
- **Concurrencia**: Uso de goroutines para manejo asíncrono
- **Limpieza automática**: Recursos se liberan automáticamente
- **Tolerancia a fallos**: Manejo robusto de errores de conexión

---

*Esta arquitectura proporciona una base sólida para comunicación en tiempo real escalable y mantenible.*
