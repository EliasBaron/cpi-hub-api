# Sistema de Notificaciones en Tiempo Real - Guía para Frontend

## 📋 Resumen

Se ha implementado un sistema de notificaciones en tiempo real para CPI Hub que notifica a los usuarios cuando alguien reacciona (like/dislike) a sus posts o comentarios. El sistema utiliza WebSockets para entregar notificaciones en tiempo real y MongoDB para persistencia.

---

## 🔄 Flujo de Notificaciones

### ¿Cuándo se emite una notificación?

Una notificación se crea automáticamente cuando:

1. Un usuario **A** reacciona (like/dislike) a un post o comentario creado por el usuario **B**
2. **Importante**: Solo se notifica si el usuario que reacciona **NO es el dueño** del contenido
3. La notificación se persiste en MongoDB y se envía por WebSocket si el usuario está conectado

### Orden de Operaciones

```
1. Usuario A reacciona → POST /v1/reactions
2. Backend persiste la reacción en MongoDB ✅
3. Backend obtiene el owner del post/comment
4. Si owner ≠ usuario que reacciona:
   → Se crea notificación en MongoDB ✅
   → Se emite evento WebSocket al owner (si está conectado)
5. Frontend del owner recibe la notificación en tiempo real
```

---

## 🔌 Endpoint WebSocket

### Conexión

**Endpoint**: `GET /v1/ws/notifications`

**Query Parameters**:
- `user_id` (requerido): ID del usuario que se conecta para recibir notificaciones

**Ejemplo de conexión**:
```javascript
const wsUrl = `ws://localhost:8080/v1/ws/notifications?user_id=${currentUser.id}`;
const socket = new WebSocket(wsUrl);
```

**Nota**: El protocolo debe ser `ws://` en desarrollo y `wss://` en producción.

---

## 📨 Formato de Mensajes

### Mensaje Recibido por WebSocket

Cuando se recibe una notificación, el mensaje tiene el siguiente formato:

```typescript
interface NotificationMessage {
  type: "notification";           // Siempre "notification"
  data: Notification;             // Datos de la notificación
  timestamp: string;               // ISO 8601 timestamp del mensaje
}

interface Notification {
  id: string;                      // ObjectID de MongoDB (hex string)
  type: "reaction";                // Tipo de notificación (actualmente solo "reaction")
  entity_type: "post" | "comment"; // Tipo de entidad a la que se reaccionó
  entity_id: number;               // ID del post o comment
  user_id: number;                 // ID del usuario que RECIBE la notificación (owner)
  read: boolean;                   // Si la notificación fue leída
  created_at: string;              // ISO 8601 timestamp de creación
}
```

### Ejemplo de Mensaje Real

```json
{
  "type": "notification",
  "data": {
    "id": "507f1f77bcf86cd799439011",
    "type": "reaction",
    "entity_type": "post",
    "entity_id": 123,
    "user_id": 456,
    "read": false,
    "created_at": "2024-01-15T10:30:00Z"
  },
  "timestamp": "2024-01-15T10:30:00Z"
}
```

**⚠️ Importante**: 
- El `user_id` en la notificación es el ID del usuario que **recibe** la notificación (el owner del contenido)
- **NO** se incluye el ID del usuario que reaccionó (anonimato de reacciones)
- El campo `created_at` es cuando se creó la notificación
- El campo `timestamp` en el mensaje es cuando se envió el mensaje WebSocket (puede ser el mismo que `created_at`)

---

## 🛠️ Implementación Frontend

### 1. Hook useNotifications

Se recomienda crear un hook personalizado para manejar las notificaciones:

```typescript
// hooks/useNotifications.ts
import { useState, useEffect, useRef } from 'react';

interface Notification {
  id: string;
  type: 'reaction';
  entity_type: 'post' | 'comment';
  entity_id: number;
  user_id: number;
  read: boolean;
  created_at: string;
}

interface NotificationMessage {
  type: 'notification';
  data: Notification;
  timestamp: string;
}

interface UseNotificationsProps {
  currentUser: { id: number } | null;
  enabled?: boolean;
}

export const useNotifications = ({ 
  currentUser, 
  enabled = true 
}: UseNotificationsProps) => {
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [unreadCount, setUnreadCount] = useState(0);
  const [connectionStatus, setConnectionStatus] = useState<
    'connecting' | 'connected' | 'disconnected' | 'error'
  >('disconnected');
  
  const socketRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);

  useEffect(() => {
    if (!enabled || !currentUser) {
      return;
    }

    const connect = () => {
      try {
        setConnectionStatus('connecting');
        
        const wsUrl = `${WEBSOCKET_BASE_URL}/v1/ws/notifications?user_id=${currentUser.id}`;
        const socket = new WebSocket(wsUrl);

        socket.onopen = () => {
          console.log('Notifications WebSocket connected');
          setConnectionStatus('connected');
          
          // Limpiar timeout de reconexión si existe
          if (reconnectTimeoutRef.current) {
            clearTimeout(reconnectTimeoutRef.current);
            reconnectTimeoutRef.current = null;
          }
        };

        socket.onmessage = (event) => {
          try {
            const message: NotificationMessage = JSON.parse(event.data);
            
            if (message.type === 'notification') {
              // Agregar nueva notificación al inicio de la lista
              setNotifications(prev => [message.data, ...prev]);
              
              // Incrementar contador de no leídas
              if (!message.data.read) {
                setUnreadCount(prev => prev + 1);
              }
              
              // Opcional: Mostrar notificación toast
              // showNotificationToast(message.data);
            }
          } catch (error) {
            console.error('Error parsing notification message:', error);
          }
        };

        socket.onerror = (error) => {
          console.error('WebSocket error:', error);
          setConnectionStatus('error');
        };

        socket.onclose = () => {
          console.log('Notifications WebSocket disconnected');
          setConnectionStatus('disconnected');
          
          // Reconexión automática después de 3 segundos
          if (enabled && currentUser) {
            reconnectTimeoutRef.current = setTimeout(() => {
              connect();
            }, 3000);
          }
        };

        socketRef.current = socket;
      } catch (error) {
        console.error('Error creating WebSocket:', error);
        setConnectionStatus('error');
      }
    };

    connect();

    // Cleanup
    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }
      if (socketRef.current) {
        socketRef.current.close();
        socketRef.current = null;
      }
    };
  }, [currentUser?.id, enabled]);

  const markAsRead = (notificationId: string) => {
    setNotifications(prev =>
      prev.map(notif =>
        notif.id === notificationId ? { ...notif, read: true } : notif
      )
    );
    setUnreadCount(prev => Math.max(0, prev - 1));
  };

  const markAllAsRead = () => {
    setNotifications(prev =>
      prev.map(notif => ({ ...notif, read: true }))
    );
    setUnreadCount(0);
  };

  return {
    notifications,
    unreadCount,
    connectionStatus,
    markAsRead,
    markAllAsRead,
  };
};
```

### 2. Factory Function para WebSocket

```typescript
// api/websocket.ts
const WEBSOCKET_BASE_URL = process.env.REACT_APP_WS_URL || 'ws://localhost:8080';

export const createNotificationWebSocket = (userId: number): WebSocket => {
  const wsUrl = `${WEBSOCKET_BASE_URL}/v1/ws/notifications?user_id=${userId}`;
  return new WebSocket(wsUrl);
};
```

### 3. Tipos TypeScript

```typescript
// types/notification.ts
export type NotificationType = 'reaction';
export type EntityType = 'post' | 'comment';

export interface Notification {
  id: string;
  type: NotificationType;
  entity_type: EntityType;
  entity_id: number;
  user_id: number;
  read: boolean;
  created_at: string;
}

export interface NotificationMessage {
  type: 'notification';
  data: Notification;
  timestamp: string;
}
```

### 4. Uso en Componentes

```typescript
// App.tsx o componente principal
import { useNotifications } from './hooks/useNotifications';

function App() {
  const { currentUser } = useAuth();
  
  const {
    notifications,
    unreadCount,
    connectionStatus,
    markAsRead,
    markAllAsRead,
  } = useNotifications({ 
    currentUser,
    enabled: !!currentUser 
  });

  return (
    <div>
      {/* Badge de notificaciones no leídas */}
      <NotificationBell unreadCount={unreadCount} />
      
      {/* Lista de notificaciones */}
      <NotificationList 
        notifications={notifications}
        onMarkAsRead={markAsRead}
        onMarkAllAsRead={markAllAsRead}
      />
      
      {/* Indicador de estado de conexión (opcional) */}
      {connectionStatus === 'connecting' && (
        <div>Conectando a notificaciones...</div>
      )}
    </div>
  );
}
```

---

## 📡 Endpoints HTTP (Futuros - Opcionales)

Aunque las notificaciones se reciben en tiempo real por WebSocket, el backend está preparado para exponer endpoints HTTP para:

- **GET** `/v1/notifications?limit=20&offset=0` - Obtener historial de notificaciones
- **GET** `/v1/notifications/unread-count` - Obtener contador de no leídas
- **PUT** `/v1/notifications/:notification_id/read` - Marcar una notificación como leída
- **PUT** `/v1/notifications/read-all` - Marcar todas como leídas

**Nota**: Estos endpoints aún no están implementados en el handler, pero el backend tiene la lógica lista. Se pueden agregar cuando sea necesario.

---

## ⚠️ Consideraciones Importantes

### 1. Reconexión Automática
- El WebSocket puede desconectarse por diversas razones (red, servidor, etc.)
- Implementar lógica de reconexión automática con backoff exponencial
- Mostrar indicador visual del estado de conexión

### 2. Manejo de Errores
- Validar que el mensaje recibido tenga el formato correcto
- Manejar errores de conexión gracefully
- No bloquear la UI si falla la conexión WebSocket

### 3. Persistencia Local (Opcional)
- Considerar guardar notificaciones en localStorage o IndexedDB
- Sincronizar con el backend cuando se recupere la conexión
- Evitar duplicados al recibir notificaciones

### 4. Notificaciones Duplicadas
- El backend puede enviar la misma notificación si hay reconexiones
- Implementar deduplicación por `id` de notificación
- Verificar si una notificación ya existe antes de agregarla

### 5. Anonimato de Reacciones
- **NO** se envía el `user_id` del usuario que reaccionó
- Solo se sabe que alguien reaccionó, pero no quién
- El frontend debe mostrar mensajes genéricos como "Alguien reaccionó a tu post"

### 6. Ping/Pong
- El backend envía pings periódicos para mantener la conexión viva
- El frontend debe responder automáticamente (el navegador lo hace por defecto)
- No es necesario implementar lógica adicional

### 7. Campo `created_at` vs `timestamp`
- `created_at`: Fecha de creación de la notificación en el backend
- `timestamp`: Fecha de envío del mensaje WebSocket
- Ambos pueden ser iguales, pero `timestamp` puede ser ligeramente posterior si hay delay en el envío

---

## 🧪 Testing

### Pruebas Manuales

1. **Conectar WebSocket**:
   ```javascript
   const ws = new WebSocket('ws://localhost:8080/v1/ws/notifications?user_id=1');
   ws.onmessage = (event) => console.log('Notification:', JSON.parse(event.data));
   ```

2. **Crear una reacción** (desde otro usuario):
   ```bash
   curl -X POST http://localhost:8080/v1/reactions \
     -H "Content-Type: application/json" \
     -d '{
       "user_id": 2,
       "entity_type": "post",
       "entity_id": 123,
       "action": "like"
     }'
   ```

3. **Verificar** que el usuario 1 (owner del post) recibe la notificación por WebSocket

---

## 📝 Ejemplo Completo de Integración

```typescript
// Ejemplo completo de componente de notificaciones
import React, { useEffect } from 'react';
import { useNotifications } from '../hooks/useNotifications';
import { useAuth } from '../contexts/AuthContext';

export const NotificationCenter: React.FC = () => {
  const { currentUser } = useAuth();
  const {
    notifications,
    unreadCount,
    connectionStatus,
    markAsRead,
    markAllAsRead,
  } = useNotifications({ currentUser });

  const handleNotificationClick = (notification: Notification) => {
    // Marcar como leída
    markAsRead(notification.id);
    
    // Navegar a la entidad
    if (notification.entity_type === 'post') {
      navigate(`/posts/${notification.entity_id}`);
    } else {
      navigate(`/comments/${notification.entity_id}`);
    }
  };

  return (
    <div className="notification-center">
      <div className="notification-header">
        <h2>Notificaciones</h2>
        {unreadCount > 0 && (
          <button onClick={markAllAsRead}>
            Marcar todas como leídas ({unreadCount})
          </button>
        )}
        <div className={`connection-status ${connectionStatus}`}>
          {connectionStatus === 'connected' && '🟢 Conectado'}
          {connectionStatus === 'connecting' && '🟡 Conectando...'}
          {connectionStatus === 'disconnected' && '🔴 Desconectado'}
        </div>
      </div>
      
      <div className="notification-list">
        {notifications.length === 0 ? (
          <p>No hay notificaciones</p>
        ) : (
          notifications.map(notification => (
            <div
              key={notification.id}
              className={`notification-item ${notification.read ? 'read' : 'unread'}`}
              onClick={() => handleNotificationClick(notification)}
            >
              <div className="notification-content">
                {notification.type === 'reaction' && (
                  <p>
                    Alguien reaccionó a tu {notification.entity_type}
                  </p>
                )}
                <small>{new Date(notification.created_at).toLocaleString()}</small>
              </div>
              {!notification.read && <div className="unread-indicator" />}
            </div>
          ))
        )}
      </div>
    </div>
  );
};
```

---

## 🔗 Recursos Adicionales

- **WebSocket API**: https://developer.mozilla.org/en-US/docs/Web/API/WebSocket
- **React Hooks**: https://react.dev/reference/react/hooks
- **TypeScript**: https://www.typescriptlang.org/docs/

---

## 📞 Soporte

Si tienes preguntas o encuentras problemas, contacta al equipo de backend.

**Última actualización**: Enero 2024

