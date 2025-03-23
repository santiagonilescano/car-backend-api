# Car Service - Microservicio de Gestión de Vehículos

Este microservicio está diseñado para gestionar el ciclo de vida de vehículos, incluyendo su historial de servicios y mantenimiento. Está construido siguiendo los principios de Domain-Driven Design (DDD).

## Estructura del Proyecto

```
.
├── cmd/
│   └── api/              # Punto de entrada de la aplicación
├── internal/
│   ├── domain/          # Entidades y reglas de negocio
│   │   ├── entities/
│   │   └── repositories/
│   ├── application/     # Casos de uso
│   └── infrastructure/  # Implementaciones concretas
│       ├── persistence/
│       └── api/
└── pkg/                 # Código compartido y utilidades
```

## Requisitos

- Go 1.21 o superior
- PostgreSQL 12 o superior

## Configuración

1. Crear una base de datos PostgreSQL:
```sql
CREATE DATABASE car_service;
```

2. Copiar el archivo de ejemplo de variables de entorno:
```bash
cp .env.example .env
```

3. Ajustar las variables en el archivo `.env` según tu configuración local.

## Ejecución

```bash
# Instalar dependencias
go mod download

# Ejecutar la aplicación
go run cmd/api/main.go
```

## API Endpoints

- `GET /health`: Verificar el estado del servicio
- Más endpoints serán documentados próximamente

## Modelo de Dominio

### Entidades Principales

1. Brand (Marca)
   - Información de la marca del vehículo
   - Relación con sus modelos

2. Model (Modelo)
   - Información del modelo
   - Relación con la marca y vehículos

3. Car (Vehículo)
   - Información básica del vehículo
   - Relación con el modelo y propietario

4. Owner (Propietario)
   - Información del propietario
   - Relación con sus vehículos

## Desarrollo

El proyecto sigue una arquitectura limpia basada en DDD con las siguientes capas:

1. Domain: Entidades y reglas de negocio core
2. Application: Casos de uso y lógica de aplicación
3. Infrastructure: Implementaciones técnicas (base de datos, API, etc.)

## Licencia

[MIT License](LICENSE) 