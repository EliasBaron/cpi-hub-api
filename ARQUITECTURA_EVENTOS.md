# Análisis Arquitectónico: EventEmitter en Use Cases

## Problema Identificado

Actualmente, el `EventEmitter` está siendo inyectado directamente en los use cases (como `PostUseCase` y `ReactionUseCase`), lo que genera una mezcla de responsabilidades:

- **Use Cases** están emitiendo eventos de infraestructura (WebSocket) después de realizar operaciones de negocio
- Se está acoplando la lógica de negocio con la infraestructura de comunicación en tiempo real
- Cada vez que un use case necesita emitir un evento, debe recibir el `EventEmitter` como dependencia

### Ejemplo Actual

```go
// post_usecase.go
type postUseCase struct {
    postRepository      domain.PostRepository
    spaceRepository     domain.SpaceRepository
    // ... otras dependencias
    eventEmitter        events.EventEmitter  // ⚠️ Infraestructura en use case
}

func (p *postUseCase) AddComment(...) {
    // Lógica de negocio
    if err := p.commentRepository.Create(ctx, comment); err != nil {
        return nil, err
    }
    
    // Emisión de evento (infraestructura)
    if p.eventEmitter != nil {
        event := &domain.Event{...}
        err := p.eventEmitter.EmitEvent(ctx, event)
    }
}
```

## ¿Es Responsabilidad de los Use Cases Emitir Eventos?

**Respuesta corta: No, no debería serlo directamente.**

### Principios Violados

1. **Separación de Responsabilidades (SRP)**: Los use cases deberían enfocarse únicamente en la lógica de negocio y orquestación, no en cómo se comunican los eventos.

2. **Inversión de Dependencias (DIP)**: Los use cases (capa de aplicación) no deberían depender directamente de implementaciones de infraestructura como WebSocket.

3. **Clean Architecture**: La capa de aplicación (use cases) no debería conocer detalles de infraestructura.

## Patrones de Diseño Alternativos

### 1. Domain Events Pattern (Recomendado)

El patrón **Domain Events** permite que las entidades del dominio emitan eventos cuando ocurren cambios significativos, y estos eventos son procesados de forma asíncrona por handlers externos.

#### Ventajas:
- ✅ Desacopla completamente los use cases de la infraestructura
- ✅ Permite múltiples handlers para un mismo evento
- ✅ Facilita testing (mock de event bus)
- ✅ Sigue principios de DDD (Domain-Driven Design)
- ✅ Permite agregar nuevos listeners sin modificar use cases

#### Estructura Propuesta:

```
Domain Layer:
  - domain/events.go (definición de eventos del dominio)
    - CommentCreatedEvent
    - ReactionCreatedEvent
    - CommentReplyCreatedEvent

Application Layer:
  - usecase/post/post_usecase.go (solo lógica de negocio)
  - usecase/reaction/reaction_usecase.go (solo lógica de negocio)

Infrastructure Layer:
  - adapters/events/domain_event_bus.go (implementación del bus)
  - adapters/events/handlers/
    - comment_event_handler.go (escucha CommentCreatedEvent)
    - reaction_event_handler.go (escucha ReactionCreatedEvent)
```

#### Flujo:
1. Use case ejecuta lógica de negocio
2. Use case publica un Domain Event en un Event Bus
3. Event Handlers (en infraestructura) escuchan los eventos
4. Handlers emiten eventos de WebSocket, guardan notificaciones, etc.

---

### 2. Mediator Pattern

El patrón **Mediator** centraliza la comunicación entre componentes a través de un mediador, evitando acoplamiento directo.

#### Ventajas:
- ✅ Desacopla emisores de receptores
- ✅ Centraliza la lógica de comunicación
- ✅ Fácil de extender con nuevos handlers

#### Estructura:
```
Application Layer:
  - usecase/post/post_usecase.go
    - Usa: domain.EventMediator

Infrastructure Layer:
  - adapters/events/event_mediator.go
    - Registra handlers para cada tipo de evento
    - Distribuye eventos a handlers apropiados
```

---

### 3. Observer Pattern

Similar al Domain Events, pero más simple. Los use cases notifican a observadores registrados.

#### Ventajas:
- ✅ Simple de implementar
- ✅ Desacopla emisores de receptores

#### Desventajas:
- ⚠️ Menos flexible que Domain Events
- ⚠️ Puede generar acoplamiento si no se diseña bien

---

### 4. Decorator Pattern (Wrapper)

Envolver los use cases con decoradores que agregan la funcionalidad de emisión de eventos.

#### Ventajas:
- ✅ No modifica código existente de use cases
- ✅ Separación clara de responsabilidades
- ✅ Fácil de activar/desactivar

#### Estructura:
```
Application Layer:
  - usecase/post/post_usecase.go (sin eventos)

Infrastructure Layer:
  - adapters/events/decorators/
    - post_usecase_event_decorator.go
      - Envuelve PostUseCase
      - Intercepta métodos
      - Emite eventos después de operaciones
```

---

## Recomendación: Domain Events Pattern

Para tu caso específico, recomiendo implementar el patrón **Domain Events** porque:

1. **Escalabilidad**: Fácil agregar nuevos listeners (email, push notifications, analytics, etc.)
2. **Testabilidad**: Los use cases se prueban sin necesidad de mockear WebSocket
3. **Mantenibilidad**: Cambios en infraestructura no afectan use cases
4. **Claridad**: La lógica de negocio queda limpia y enfocada

### Implementación Sugerida

#### 1. Definir Domain Events

```go
// domain/events/domain_events.go
package events

type DomainEvent interface {
    EventType() string
    OccurredAt() time.Time
}

type CommentCreatedEvent struct {
    CommentID    int
    PostID       int
    CreatedBy    int
    PostOwnerID  int
    Content      string
    IsReply      bool
    ParentID     *int
    OccurredAt   time.Time
}

func (e CommentCreatedEvent) EventType() string {
    return "comment_created"
}

type ReactionCreatedEvent struct {
    ReactionID   string
    UserID       int
    OwnerUserID  int
    EntityType   string
    EntityID     int
    Action       string
    PostID       *int
    OccurredAt   time.Time
}
```

#### 2. Event Bus Interface (Domain)

```go
// domain/events/event_bus.go
package events

type EventBus interface {
    Publish(ctx context.Context, event DomainEvent) error
    Subscribe(eventType string, handler EventHandler)
}

type EventHandler interface {
    Handle(ctx context.Context, event DomainEvent) error
}
```

#### 3. Use Cases (Solo Lógica de Negocio)

```go
// usecase/post/post_usecase.go
type postUseCase struct {
    postRepository    domain.PostRepository
    commentRepository domain.CommentRepository
    eventBus          domain.EventBus  // Solo interfaz del dominio
}

func (p *postUseCase) AddComment(ctx context.Context, commentDTO dto.CreateComment) {
    // Lógica de negocio
    comment := commentDTO.ToDomain()
    // ... validaciones y creación
    
    // Publicar evento del dominio (no infraestructura)
    if post.CreatedBy != comment.CreatedBy {
        event := events.CommentCreatedEvent{
            CommentID:   comment.ID,
            PostID:      comment.PostID,
            CreatedBy:   comment.CreatedBy,
            PostOwnerID: post.CreatedBy,
            Content:     comment.Content,
            OccurredAt:  time.Now(),
        }
        p.eventBus.Publish(ctx, event)
    }
    
    return comment
}
```

#### 4. Event Handlers (Infraestructura)

```go
// infrastructure/adapters/events/handlers/comment_handler.go
type CommentEventHandler struct {
    eventEmitter events.EventEmitter
}

func (h *CommentEventHandler) Handle(ctx context.Context, event events.DomainEvent) error {
    commentEvent, ok := event.(events.CommentCreatedEvent)
    if !ok {
        return nil
    }
    
    // Convertir Domain Event a Event de infraestructura
    wsEvent := &domain.Event{
        Type:         "comment_created",
        UserID:       commentEvent.CreatedBy,
        TargetUserID: commentEvent.PostOwnerID,
        Metadata: map[string]interface{}{
            "post_id":    commentEvent.PostID,
            "comment_id": commentEvent.CommentID,
        },
        Timestamp: commentEvent.OccurredAt,
    }
    
    return h.eventEmitter.EmitEvent(ctx, wsEvent)
}
```

#### 5. Event Bus Implementation (Infraestructura)

```go
// infrastructure/adapters/events/in_memory_event_bus.go
type InMemoryEventBus struct {
    handlers map[string][]domain.EventHandler
    mutex    sync.RWMutex
}

func (b *InMemoryEventBus) Publish(ctx context.Context, event domain.DomainEvent) error {
    b.mutex.RLock()
    handlers := b.handlers[event.EventType()]
    b.mutex.RUnlock()
    
    for _, handler := range handlers {
        go handler.Handle(ctx, event) // Asíncrono
    }
    return nil
}
```

---

## Comparación de Patrones

| Patrón | Complejidad | Desacoplamiento | Escalabilidad | Testabilidad |
|--------|-------------|-----------------|---------------|---------------|
| **Domain Events** | Media-Alta | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| **Mediator** | Media | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| **Observer** | Baja | ⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐ |
| **Decorator** | Baja | ⭐⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐ |
| **Actual (EventEmitter directo)** | Baja | ⭐ | ⭐⭐ | ⭐⭐ |

---

## Conclusión

**No, no es responsabilidad de los use cases emitir eventos directamente.** Los use cases deben enfocarse en la lógica de negocio y orquestación. La emisión de eventos debe ser una consecuencia de las acciones del dominio, manejada por la capa de infraestructura.

El patrón **Domain Events** es la mejor opción para tu arquitectura porque:
- Mantiene los use cases limpios y enfocados
- Permite múltiples consumidores de eventos sin modificar use cases
- Facilita testing y mantenimiento
- Sigue principios de Clean Architecture y DDD

---

## Referencias

- [Domain Events Pattern - Martin Fowler](https://martinfowler.com/eaaDev/DomainEvent.html)
- [Clean Architecture - Robert C. Martin](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design - Eric Evans](https://www.domainlanguage.com/ddd/)

