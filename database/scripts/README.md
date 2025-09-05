# CPI Hub API - Database Scripts

Este directorio contiene los scripts de base de datos para la aplicación CPI Hub API.

## Estructura de archivos

- `01_ddl.sql` - Script DDL para crear todas las tablas y estructuras de la base de datos
- `02_sample_data.sql` - Script con datos de prueba para poblar la base de datos
- `README.md` - Este archivo de documentación

## Uso

### 1. Crear la base de datos

Primero, asegúrate de tener PostgreSQL instalado y crear una base de datos:

```sql
CREATE DATABASE cpi_hub;
```

### 2. Ejecutar el DDL

Ejecuta el script DDL para crear todas las tablas:

```bash
psql -d cpi_hub -f database/scripts/01_ddl.sql
```

### 3. Cargar datos de prueba

Ejecuta el script de datos de muestra:

```bash
psql -d cpi_hub -f database/scripts/02_sample_data.sql
```

## Estructura de la base de datos

### Tablas principales

- **users** - Usuarios del sistema
- **spaces** - Espacios de trabajo/comunidades
- **posts** - Publicaciones dentro de los espacios
- **comments** - Comentarios en las publicaciones
- **user_spaces** - Relación many-to-many entre usuarios y espacios

### Datos de prueba incluidos

- 8 usuarios de ejemplo con diferentes roles
- 6 espacios temáticos (Desarrollo Web, DevOps, ML, etc.)
- Relaciones usuario-espacio configuradas
- 6 publicaciones de ejemplo
- 12 comentarios de ejemplo

## Notas

- Las contraseñas en los datos de prueba están hasheadas con bcrypt
- Se incluyen triggers para actualizar automáticamente el campo `updated_at`
- Se crean índices para optimizar las consultas más comunes
- Las imágenes de perfil usan placeholders de ejemplo

## Personalización

Puedes modificar el archivo `02_sample_data.sql` para agregar más datos de prueba según tus necesidades. Los usuarios de ejemplo incluyen:

1. Juan Pérez - Desarrollador Full Stack
2. María González - DevOps Engineer
3. Carlos Rodríguez - Data Scientist
4. Ana Martín - Mobile Developer
5. Luis Fernández - Backend Developer
6. Laura López - Frontend Developer
7. Pedro Sánchez - ML Engineer
8. Sofia Ramírez - Full Stack Developer
