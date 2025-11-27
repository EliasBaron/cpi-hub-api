# Guía Completa de Despliegue - CPI Hub API

## 📋 Tabla de Contenidos

1. [Análisis de la Aplicación](#análisis-de-la-aplicación)
2. [Plan Gratuito y Límites](#plan-gratuito-y-límites)
3. [Configuración de MongoDB Atlas](#configuración-de-mongodb-atlas)
4. [Configuración de PostgreSQL en Render](#configuración-de-postgresql-en-render)
5. [Modificaciones Necesarias en el Código](#modificaciones-necesarias-en-el-código)
6. [Despliegue en Render](#despliegue-en-render)
7. [Gestión de Secrets y Variables de Entorno](#gestión-de-secrets-y-variables-de-entorno)
8. [Verificación y Testing](#verificación-y-testing)
9. [Consideraciones de Seguridad](#consideraciones-de-seguridad)
10. [Troubleshooting](#troubleshooting)

---

## 🔍 Análisis de la Aplicación

### Stack Tecnológico
- **Framework**: Gin (Go)
- **Base de Datos Relacional**: PostgreSQL (usuarios, posts, comentarios, espacios, mensajes)
- **Base de Datos NoSQL**: MongoDB (reacciones, notificaciones, news)
- **WebSockets**: Para notificaciones en tiempo real
- **CORS**: Configurado para `localhost:3000` (necesita actualización)

### Configuraciones Actuales
- **PostgreSQL**: Hardcodeado en `clients.go` (localhost:5432)
- **MongoDB**: Hardcodeado en `clients.go` (localhost:27017)
- **CORS**: Solo permite `http://localhost:3000`
- **Puerto**: Lee de variable de entorno `PORT` (default: 8080)

### Cambios Necesarios
1. ✅ Leer conexiones de bases de datos desde variables de entorno
2. ✅ Configurar CORS para producción (permitir dominio del frontend)
3. ✅ Asegurar que el puerto se lea correctamente (ya está implementado)

---

## 💰 Plan Gratuito y Límites

### Render Free Tier

**Servicios Web:**
- ✅ **750 horas gratuitas por mes** (suficiente para 1 servicio 24/7)
- ✅ Auto-deploy desde GitHub
- ✅ HTTPS automático
- ✅ Logs en tiempo real
- ⚠️ **Servicios se duermen después de 15 minutos de inactividad** (se despiertan automáticamente en la primera petición)

**PostgreSQL:**
- ✅ **1 GB de almacenamiento**
- ⚠️ **Expira después de 30 días** (anteriormente 90 días)
- ✅ Backup automático
- ⚠️ **Sin tarjeta de crédito requerida**

**Créditos:**
- Render no usa un sistema de créditos tradicional
- Los servicios gratuitos tienen límites de tiempo y recursos
- Si necesitas más recursos, puedes actualizar a un plan de pago

**Referencias:**
- [Render Free Tier](https://render.com/docs/free)
- [PostgreSQL Free Tier Changes](https://render.com/changelog/free-postgresql-instances-now-expire-after-30-days-previously-90)

### MongoDB Atlas Free Tier (M0)

**Características:**
- ✅ **512 MB de almacenamiento**
- ✅ **Clúster compartido** (no dedicado)
- ✅ **Sin límite de tiempo** (permanente mientras uses el clúster)
- ✅ **Alta disponibilidad** (replicación automática)
- ✅ **Backup automático** (últimas 2 semanas)
- ✅ **Sin tarjeta de crédito requerida**
- ⚠️ **Límite de conexiones**: 500 conexiones simultáneas
- ⚠️ **Performance**: Compartido con otros usuarios (puede ser más lento)

**Referencias:**
- [MongoDB Atlas Free Tier](https://www.mongodb.com/cloud/atlas/pricing)

---

## 🍃 Configuración de MongoDB Atlas

### Paso 1: Crear Cuenta y Proyecto

1. Visita [MongoDB Atlas](https://www.mongodb.com/cloud/atlas/register)
2. Crea una cuenta gratuita (puedes usar Google, GitHub, o email)
3. Una vez dentro, crea un nuevo **Project**:
   - Click en "New Project"
   - Nombre: `CPI Hub` (o el que prefieras)
   - Click "Create Project"

### Paso 2: Crear Cluster Gratuito

1. En el dashboard, click en **"Build a Database"**
2. Selecciona el plan **"M0 (Free)"**
3. **Configuración del Cluster:**
   - **Cloud Provider**: AWS, Google Cloud, o Azure (elige el más cercano a tu región)
   - **Region**: Selecciona la región más cercana (ej: `us-east-1` para AWS)
   - **Cluster Name**: `cpi-hub-cluster` (o el que prefieras)
4. Click **"Create Cluster"**
5. ⏳ Espera 3-5 minutos mientras se crea el cluster

### Paso 3: Configurar Acceso a la Base de Datos

#### 3.1. Crear Usuario de Base de Datos

1. En el dashboard, ve a **"Database Access"** (menú lateral izquierdo)
2. Click **"Add New Database User"**
3. **Método de Autenticación**: "Password"
4. **Username**: `cpi-hub-user` (o el que prefieras)
5. **Password**: Genera una contraseña segura (guárdala, la necesitarás)
6. **Database User Privileges**: "Atlas admin" (o "Read and write to any database")
7. Click **"Add User"**

#### 3.2. Configurar Network Access (IP Whitelist)

1. Ve a **"Network Access"** (menú lateral izquierdo)
2. Click **"Add IP Address"**
3. Para desarrollo/testing, puedes usar:
   - **"Add Current IP Address"** (agrega tu IP actual)
   - **"Allow Access from Anywhere"** (0.0.0.0/0) - ⚠️ **Solo para desarrollo, no recomendado para producción**
4. Click **"Confirm"**

> **Nota para Producción**: Render tiene IPs dinámicas. Para producción, necesitarás permitir acceso desde cualquier IP (0.0.0.0/0) o configurar un Private Endpoint (requiere plan de pago).

### Paso 4: Obtener Connection String

1. En el dashboard, click en **"Connect"** (botón en tu cluster)
2. Selecciona **"Connect your application"**
3. **Driver**: "Go" (versión 1.18 o superior)
4. Copia la **Connection String**, se verá así:
   ```
   mongodb+srv://<username>:<password>@cpi-hub-cluster.xxxxx.mongodb.net/?retryWrites=true&w=majority
   ```
5. Reemplaza `<username>` y `<password>` con las credenciales que creaste
6. Agrega el nombre de la base de datos al final:
   ```
   mongodb+srv://cpi-hub-user:TU_PASSWORD@cpi-hub-cluster.xxxxx.mongodb.net/cpihub?retryWrites=true&w=majority
   ```

### Paso 5: Verificar Conexión (Opcional)

Puedes probar la conexión usando MongoDB Compass o el MongoDB Shell:

```bash
# Instalar MongoDB Shell (opcional)
# macOS: brew install mongosh
# Luego conectar:
mongosh "mongodb+srv://cpi-hub-user:TU_PASSWORD@cpi-hub-cluster.xxxxx.mongodb.net/cpihub"
```

---

## 🐘 Configuración de PostgreSQL en Render

### Paso 1: Crear Base de Datos PostgreSQL

1. Inicia sesión en [Render Dashboard](https://dashboard.render.com)
2. Click en **"New +"** (esquina superior derecha)
3. Selecciona **"PostgreSQL"**
4. **Configuración:**
   - **Name**: `cpi-hub-db` (o el que prefieras)
   - **Database**: `cpihub` (o el que prefieras)
   - **User**: Se genera automáticamente
   - **Region**: Selecciona la misma región que usarás para tu servicio web
   - **PostgreSQL Version**: Deja la última versión (15 o superior)
   - **Plan**: **Free** (1 GB storage, expira en 30 días)
5. Click **"Create Database"**
6. ⏳ Espera 2-3 minutos mientras se crea la base de datos

### Paso 2: Obtener Credenciales de Conexión

1. Una vez creada, ve a la página de tu base de datos
2. En la sección **"Connections"**, encontrarás:
   - **Host**: `dpg-xxxxx-a.oregon-postgres.render.com`
   - **Port**: `5432`
   - **Database**: `cpihub`
   - **User**: `cpi_hub_db_user` (o similar)
   - **Password**: Se muestra una vez (cópiala inmediatamente)
   - **Internal Database URL**: `postgres://user:password@host:port/database`
   - **External Database URL**: Similar pero para conexiones externas

3. **Guarda estas credenciales**, especialmente la **Internal Database URL** que usarás en Render

### Paso 3: Ejecutar Migraciones (Opcional - Manual)

Si quieres ejecutar las migraciones manualmente antes del deploy:

1. Obtén la **External Database URL** de Render
2. Conéctate usando `psql`:
   ```bash
   psql "postgres://user:password@host:port/database"
   ```
3. Ejecuta las migraciones desde `database/migrations/` o el schema se creará automáticamente al iniciar la app

> **Nota**: Tu aplicación ejecuta `schema.EnsureSchema()` al iniciar, así que las tablas se crearán automáticamente.

---

## 🔧 Modificaciones Necesarias en el Código

### 1. Actualizar `internal/app/dependencies/clients.go`

Necesitamos modificar las funciones para leer desde variables de entorno:

```go
package dependencies

import (
	"context"
	"cpi-hub-api/database/schema"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"database/sql"

	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func newMongoDBClient() (*mongo.Client, error) {
	// Leer desde variable de entorno
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017" // Fallback para desarrollo local
	}

	databaseName := os.Getenv("MONGODB_DATABASE")
	if databaseName == "" {
		databaseName = "cpihub"
	}

	timeout := 10 * time.Second
	if timeoutStr := os.Getenv("MONGODB_TIMEOUT"); timeoutStr != "" {
		if parsedTimeout, err := time.ParseDuration(timeoutStr); err == nil {
			timeout = parsedTimeout
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error conectando a MongoDB: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error verificando conexión a MongoDB: %w", err)
	}

	log.Printf("Conectado exitosamente a MongoDB")
	return client, nil
}

func GetMongoDatabase() (*mongo.Database, error) {
	client, err := newMongoDBClient()
	if err != nil {
		return nil, err
	}

	databaseName := os.Getenv("MONGODB_DATABASE")
	if databaseName == "" {
		databaseName = "cpihub"
	}

	return client.Database(databaseName), nil
}

func NewPostgreSQLClient() (*sql.DB, error) {
	// Intentar usar DATABASE_URL primero (formato completo)
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL != "" {
		db, err := sql.Open("postgres", databaseURL)
		if err != nil {
			return nil, fmt.Errorf("error opening connection to PostgreSQL: %w", err)
		}

		if err := db.Ping(); err != nil {
			return nil, fmt.Errorf("error verifying connection to PostgreSQL: %w", err)
		}

		if err = schema.EnsureSchema(db); err != nil {
			return nil, fmt.Errorf("error ensuring database schema: %w", err)
		}

		log.Printf("Successfully connected to PostgreSQL using DATABASE_URL")
		return db, nil
	}

	// Fallback: leer variables individuales
	config := PostgreSQLConfig{
		Host:     getEnv("POSTGRES_HOST", "localhost"),
		Port:     getEnvAsInt("POSTGRES_PORT", 5432),
		User:     getEnv("POSTGRES_USER", "postgres"),
		Password: getEnv("POSTGRES_PASSWORD", "rootroot"),
		Database: getEnv("POSTGRES_DB", "cpihub"),
		SSLMode:  getEnv("POSTGRES_SSLMODE", "require"),
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.Database, config.SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening connection to PostgreSQL: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error verifying connection to PostgreSQL: %w", err)
	}

	if err = schema.EnsureSchema(db); err != nil {
		return nil, fmt.Errorf("error ensuring database schema: %w", err)
	}

	log.Printf("Successfully connected to PostgreSQL at %s:%d", config.Host, config.Port)
	return db, nil
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
```

### 2. Actualizar CORS en `cmd/api/main.go`

Necesitamos permitir el dominio de producción:

```go
package main

import (
	"log"
	"os"
	"strings"

	"cpi-hub-api/internal/app/dependencies"
	"cpi-hub-api/internal/infrastructure/entrypoint/router"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	app := gin.Default()

	// Configurar CORS desde variables de entorno
	allowedOrigins := []string{"http://localhost:3000"} // Default para desarrollo
	
	if originsEnv := os.Getenv("CORS_ALLOWED_ORIGINS"); originsEnv != "" {
		allowedOrigins = strings.Split(originsEnv, ",")
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length"},
		MaxAge:           12 * time.Hour,
	}))

	router.LoadRoutes(app, dependencies.Build())

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor iniciando en el puerto %s", port)
	if err := app.Run(":" + port); err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
```

> **Nota**: Necesitarás agregar `"time"` al import si no está ya.

---

## 🚀 Despliegue en Render

### Paso 1: Preparar el Repositorio

1. Asegúrate de que todos los cambios estén commiteados y pusheados a GitHub:
   ```bash
   git add .
   git commit -m "feat: configure environment variables for deployment"
   git push origin develop
   ```

2. Verifica que el repositorio esté actualizado en GitHub

### Paso 2: Crear Web Service en Render

1. En [Render Dashboard](https://dashboard.render.com), click **"New +"**
2. Selecciona **"Web Service"**
3. **Conectar Repositorio:**
   - Si es la primera vez, autoriza Render para acceder a tu GitHub
   - Selecciona el repositorio `cpi-hub-api`
   - Selecciona la rama `develop` (o `main` según tu preferencia)

### Paso 3: Configurar el Servicio

**Configuración Básica:**
- **Name**: `cpi-hub-api` (o el que prefieras)
- **Region**: Selecciona la misma región que tu base de datos PostgreSQL
- **Branch**: `develop` (o la rama que uses)
- **Root Directory**: Dejar vacío (o `./` si está en la raíz)

**Build & Deploy:**
- **Environment**: `Go`
- **Build Command**: `go build -tags netgo -ldflags '-s -w' -o app`
- **Start Command**: `./app`

**Plan:**
- Selecciona **"Free"** (750 horas/mes)

### Paso 4: Configurar Variables de Entorno

Antes de hacer el deploy, configura las variables de entorno:

1. En la configuración del servicio, ve a la pestaña **"Environment"**
2. Agrega las siguientes variables:

   **PostgreSQL:**
   ```
   DATABASE_URL=postgres://user:password@host:port/database
   ```
   (Usa la **Internal Database URL** de tu base de datos PostgreSQL en Render)

   **MongoDB:**
   ```
   MONGODB_URI=mongodb+srv://cpi-hub-user:TU_PASSWORD@cpi-hub-cluster.xxxxx.mongodb.net/cpihub?retryWrites=true&w=majority
   MONGODB_DATABASE=cpihub
   ```

   **CORS (Opcional - para producción):**
   ```
   CORS_ALLOWED_ORIGINS=https://tu-frontend.com,https://www.tu-frontend.com
   ```

   **Puerto (Opcional - Render lo configura automáticamente):**
   ```
   PORT=10000
   ```
   (Render asigna el puerto automáticamente, pero puedes especificarlo)

3. Click **"Save Changes"**

### Paso 5: Conectar Base de Datos PostgreSQL

1. En la configuración del servicio, ve a la pestaña **"Environment"**
2. Busca la sección **"Add Environment Variable"**
3. Click en **"Link Database"** o busca tu base de datos PostgreSQL
4. Render automáticamente agregará `DATABASE_URL` con la conexión interna

> **Nota**: Si ya agregaste `DATABASE_URL` manualmente, puedes eliminarla y usar el "Link Database" que es más seguro.

### Paso 6: Iniciar el Deploy

1. Click **"Create Web Service"** o **"Save Changes"** si ya existe
2. Render comenzará a construir y desplegar tu aplicación
3. Puedes ver el progreso en la pestaña **"Logs"**
4. ⏳ El deploy puede tardar 5-10 minutos la primera vez

### Paso 7: Verificar el Deploy

1. Una vez completado, tu aplicación estará disponible en:
   ```
   https://cpi-hub-api.onrender.com
   ```
   (El nombre puede variar según el que hayas elegido)

2. Verifica los logs para asegurarte de que:
   - ✅ La conexión a PostgreSQL fue exitosa
   - ✅ La conexión a MongoDB fue exitosa
   - ✅ El servidor está escuchando en el puerto correcto

---

## 🔐 Gestión de Secrets y Variables de Entorno

### Variables de Entorno Requeridas

| Variable | Descripción | Ejemplo | Requerida |
|----------|-------------|---------|-----------|
| `DATABASE_URL` | Connection string de PostgreSQL | `postgres://user:pass@host:port/db` | ✅ Sí |
| `MONGODB_URI` | Connection string de MongoDB Atlas | `mongodb+srv://user:pass@cluster.mongodb.net/db` | ✅ Sí |
| `MONGODB_DATABASE` | Nombre de la base de datos MongoDB | `cpihub` | ⚠️ Opcional |
| `PORT` | Puerto del servidor | `10000` | ⚠️ Opcional (Render lo asigna) |
| `CORS_ALLOWED_ORIGINS` | Orígenes permitidos (separados por coma) | `https://app.com,https://www.app.com` | ⚠️ Opcional |

### Mejores Prácticas

1. **Nunca commitees secrets al repositorio**
   - Usa `.gitignore` para archivos con secrets
   - Usa variables de entorno siempre

2. **Usa diferentes valores para desarrollo y producción**
   - Desarrollo: valores locales
   - Producción: valores desde Render/MongoDB Atlas

3. **Rota las contraseñas periódicamente**
   - MongoDB Atlas: Regenera usuarios cada 3-6 meses
   - PostgreSQL: Regenera contraseñas desde Render dashboard

4. **Usa Internal Database URL en Render**
   - Render proporciona URLs internas que son más seguras y rápidas
   - No exponen la base de datos a internet

---

## ✅ Verificación y Testing

### 1. Verificar Conexiones

Una vez desplegado, verifica que las conexiones funcionen:

```bash
# Ver logs en Render
# Dashboard > Tu Servicio > Logs

# Deberías ver:
# ✅ "Successfully connected to PostgreSQL..."
# ✅ "Conectado exitosamente a MongoDB"
# ✅ "Servidor iniciando en el puerto 10000"
```

### 2. Probar Endpoints

```bash
# Health check (si tienes uno)
curl https://cpi-hub-api.onrender.com/health

# O probar un endpoint real
curl https://cpi-hub-api.onrender.com/v1/users
```

### 3. Verificar Base de Datos

**PostgreSQL:**
- Las tablas se crean automáticamente al iniciar
- Puedes verificar en Render > PostgreSQL > "Connect" > usar `psql`

**MongoDB:**
- Conéctate usando MongoDB Compass o `mongosh`
- Verifica que las colecciones se creen cuando se usen

### 4. Probar WebSockets

Si tu aplicación usa WebSockets, verifica que funcionen:
- Render soporta WebSockets en el plan gratuito
- Prueba la conexión desde tu frontend

---

## 🛡️ Consideraciones de Seguridad

### 1. MongoDB Atlas

- ✅ **IP Whitelist**: Para producción, considera usar Private Endpoints (requiere plan de pago) o restringir IPs conocidas
- ✅ **Usuarios con permisos mínimos**: No uses "Atlas admin" en producción, crea usuarios con permisos específicos
- ✅ **Contraseñas fuertes**: Usa contraseñas generadas automáticamente
- ✅ **Habilita MFA**: Activa autenticación de dos factores en tu cuenta de Atlas

### 2. PostgreSQL en Render

- ✅ **Internal Database URL**: Usa siempre la URL interna en Render (más segura)
- ✅ **SSL Required**: Render fuerza SSL, asegúrate de usar `sslmode=require`
- ✅ **Backups**: Las bases de datos gratuitas tienen backups automáticos, pero expiran en 30 días

### 3. Aplicación

- ✅ **HTTPS**: Render proporciona HTTPS automático
- ✅ **CORS**: Configura solo los orígenes necesarios
- ✅ **Rate Limiting**: Considera agregar rate limiting para proteger tu API
- ✅ **Logs**: No loguees información sensible (contraseñas, tokens)

### 4. Variables de Entorno

- ✅ **Render Secrets**: Render encripta las variables de entorno
- ✅ **No las expongas**: Nunca las incluyas en logs públicos o código
- ✅ **Rota periódicamente**: Cambia las contraseñas cada 3-6 meses

---

## 🔧 Troubleshooting

### Problema: Error de conexión a PostgreSQL

**Síntomas:**
```
error opening connection to PostgreSQL: connection refused
```

**Soluciones:**
1. Verifica que `DATABASE_URL` esté correctamente configurada
2. Asegúrate de usar la **Internal Database URL** (no External)
3. Verifica que la base de datos esté activa en Render
4. Revisa que el formato de la URL sea correcto: `postgres://user:password@host:port/database`

### Problema: Error de conexión a MongoDB

**Síntomas:**
```
error conectando a MongoDB: connection timeout
```

**Soluciones:**
1. Verifica que `MONGODB_URI` esté correctamente configurada
2. Asegúrate de que la IP de Render esté en la whitelist de MongoDB Atlas
   - Para desarrollo: Agrega `0.0.0.0/0` (permite todas las IPs)
3. Verifica que el usuario y contraseña sean correctos
4. Asegúrate de que el cluster esté activo en MongoDB Atlas

### Problema: CORS Error

**Síntomas:**
```
Access to fetch at 'https://cpi-hub-api.onrender.com' from origin 'https://tu-frontend.com' has been blocked by CORS policy
```

**Soluciones:**
1. Agrega tu dominio frontend a `CORS_ALLOWED_ORIGINS` en Render
2. Verifica que el formato sea correcto: `https://tu-frontend.com,https://www.tu-frontend.com`
3. Reinicia el servicio después de cambiar las variables de entorno

### Problema: Servicio se duerme

**Síntomas:**
- Primera petición después de 15 minutos de inactividad tarda mucho

**Explicación:**
- Esto es normal en el plan gratuito de Render
- El servicio se "duerme" después de 15 minutos de inactividad
- Se "despierta" automáticamente en la primera petición (puede tardar 30-60 segundos)

**Soluciones:**
1. Usa un servicio de "ping" periódico (cada 10 minutos) para mantenerlo activo
2. Actualiza a un plan de pago si necesitas que esté siempre activo

### Problema: Build Fails

**Síntomas:**
```
Build failed: go: cannot find module
```

**Soluciones:**
1. Verifica que `go.mod` esté actualizado
2. Asegúrate de que todas las dependencias estén en `go.mod`
3. Verifica que la versión de Go sea compatible (Render usa Go 1.21+ por defecto)
4. Puedes especificar la versión en `go.mod`: `go 1.24`

### Problema: Puerto no disponible

**Síntomas:**
```
Error al iniciar el servidor: listen tcp :8080: bind: address already in use
```

**Soluciones:**
1. Render asigna el puerto automáticamente a través de `PORT`
2. No hardcodees el puerto, siempre usa `os.Getenv("PORT")`
3. Si Render no asigna `PORT`, contacta soporte

---

## 📝 Checklist Final

Antes de considerar el despliegue completo, verifica:

### MongoDB Atlas
- [ ] Cluster M0 creado y activo
- [ ] Usuario de base de datos creado
- [ ] IP Whitelist configurada (0.0.0.0/0 para desarrollo)
- [ ] Connection string obtenida y probada
- [ ] Contraseña guardada de forma segura

### PostgreSQL en Render
- [ ] Base de datos creada
- [ ] Internal Database URL obtenida
- [ ] Credenciales guardadas de forma segura
- [ ] Base de datos activa

### Código
- [ ] `clients.go` actualizado para leer variables de entorno
- [ ] `main.go` actualizado para CORS configurable
- [ ] Cambios commiteados y pusheados a GitHub
- [ ] Código probado localmente con variables de entorno

### Render
- [ ] Web Service creado
- [ ] Repositorio conectado
- [ ] Build command configurado
- [ ] Start command configurado
- [ ] Variables de entorno configuradas
- [ ] Base de datos PostgreSQL vinculada
- [ ] Deploy exitoso
- [ ] Logs verificados (sin errores)

### Testing
- [ ] Conexión a PostgreSQL verificada
- [ ] Conexión a MongoDB verificada
- [ ] Endpoints probados
- [ ] WebSockets funcionando (si aplica)
- [ ] CORS configurado correctamente

---

## 🎉 ¡Listo!

Tu aplicación debería estar desplegada y funcionando. Si encuentras algún problema, revisa la sección de Troubleshooting o los logs en Render.

### Próximos Pasos

1. **Monitoreo**: Configura alertas en Render y MongoDB Atlas
2. **Backups**: Considera configurar backups adicionales para producción
3. **Escalabilidad**: Si necesitas más recursos, considera actualizar a planes de pago
4. **Documentación**: Documenta tu API para que otros desarrolladores puedan usarla

### Recursos Útiles

- [Render Documentation](https://render.com/docs)
- [MongoDB Atlas Documentation](https://www.mongodb.com/docs/atlas/)
- [Go Documentation](https://go.dev/doc/)
- [Gin Framework Documentation](https://gin-gonic.com/docs/)

---

**Última actualización**: Noviembre 2024


