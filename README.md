# inventory-service-go

Implementación en **Go** del servicio de inventario. Es un módulo de este
mismo repositorio (`inventory-service-go/`) y una **alternativa** a
`inventory-service` (Spring Boot): se usa uno u otro, **nunca ambos activos
contra `inventory_db`**. Ver `../docs/decisions/ADR-002-inventory-service-go.md`.

## Responsabilidades

Igual que `inventory-service`:

* API REST de inventario.
* Lógica de negocio y persistencia en `inventory_db`.
* Validación del JWT emitido por `auth-service`.
* Autorización (roles/permisos: decisión futura, aún no implementada).

No implementa login, password hashing ni UI. `auth-service` sigue siendo el
único emisor de JWT.

## Stack

* Go (stdlib `net/http`, `database/sql`).
* PostgreSQL (`github.com/lib/pq`).
* JWT HS256 (`github.com/golang-jwt/jwt/v5`).

## Estructura

```text
cmd/api/main.go          # entrypoint
internal/config/         # variables de entorno
internal/db/             # conexión PostgreSQL
internal/auth/           # validación del JWT
internal/inventory/      # dominio, Repository e implementación Postgres
internal/api/            # handlers, router y middleware
```

## Configuración

Copiar `.env.example` a `.env` y completarlo. Todas las variables se leen del
entorno; este proyecto no incluye librería de `.env`, cargalas en tu shell o
gestor de contenedores. Nunca commitear secretos.

| Variable | Descripción |
|---|---|
| `PORT` | Puerto HTTP (default `8080`) |
| `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME` | Conexión a `inventory_db` |
| `DB_SSLMODE` | Modo SSL de PostgreSQL (default `disable` para desarrollo) |
| `JWT_SECRET` | Mismo secreto HS256 que `auth-service` (obligatorio) |

## Base de datos

Usar el esquema del repositorio:

```text
../DDBB/inventory.sql
../DDBB/insert-inventory.sql
```

No crear foreign keys hacia `auth_db`. `movements.user_id` y
`transfers.user_id` son referencias lógicas al `sub` del JWT.

## Ejecutar

```bash
go mod tidy     # genera go.sum la primera vez
make test
make run
```

En PowerShell, cargar variables por sesión:

```powershell
$env:JWT_SECRET = "..."
$env:DB_PASSWORD = "..."
go run ./cmd/api
```

## Endpoints implementados

```text
GET    /health
GET    /api/products
GET    /api/products/{id}
POST   /api/products
PUT    /api/products/{id}
DELETE /api/products/{id}
```

Todas las rutas `/api/` requieren `Authorization: Bearer <JWT>`.
El contrato completo está en `../docs/api-contract.md`.

## Pendiente

Warehouses, categories, suppliers, inventory, movements y transfers.
Los endpoints deben implementarse con el mismo contrato que `inventory-service`.

## Cambio entre implementaciones

1. Detener el servicio de inventario activo.
2. Arrancar este servicio apuntando a la misma `inventory_db`.
3. Actualizar `INVENTORY_SERVICE_URL` en `inventory-front/.env`.
4. Verificar `/health`, login y una operación de lectura/escritura.

Nunca dejar Spring Boot y Go escribiendo al mismo tiempo en `inventory_db`.
