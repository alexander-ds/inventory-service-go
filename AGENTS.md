# AGENTS.md — inventory-service-go

## Propósito

Implementación en Go de la API de inventario. Es **alternativa** a
`inventory-service` (Spring Boot): se usa uno u otro, nunca ambos activos
contra `inventory_db`. Ver `../docs/decisions/ADR-002-inventory-service-go.md`.

## Responsabilidades

* API REST (`net/http`).
* Lógica de negocio de inventario.
* Persistencia en `inventory_db`.
* Validación del JWT emitido por `auth-service`.
* Autorización (roles/permisos: decisión futura).
* Manejo de errores HTTP (401, 403, 400, 404).

## No responsabilidades

No implementar:

* Login con email/password.
* Password hashing.
* Emisión de JWT.
* Acceso a `auth_db`.
* UI.

## Arquitectura

Preferir:

```text
Handler (internal/api)
    ↓
Service (internal/inventory)
    ↓
Repository interface
    ↓
PostgresRepository (internal/db)
    ↓
inventory_db
```

Los handlers no contienen reglas de negocio. La capa `Service` depende de una
interfaz `Repository`, para poder mockearla en tests.

## Seguridad

* Toda ruta `/api/` exige `Authorization: Bearer <JWT>`.
* Validar HS256, firma y expiración con `JWT_SECRET` del entorno.
* No inventar claims: actuales `sub`, `email`, `name` (opcional), `iat`, `exp`.
* Roles y permisos son futuros; no implementarlos hasta documentar el contrato.
* Nunca hardcodear secretos ni loguear tokens completos o passwords.
* Usar queries parametrizadas (`$1`, `$2`), nunca concatenar valores del usuario.

## Base de datos

* Usar el esquema `../DDBB/inventory.sql`.
* No crear foreign keys hacia `auth_db`.
* `movements.user_id` y `transfers.user_id` son referencia lógica al claim `sub`.

## Contrato

Seguir `../docs/api-contract.md` y `../docs/authentication.md`.
Cualquier cambio de endpoints, JWT o esquema requiere actualizar la
documentación y, si aplica, un ADR.

## Testing

Con mocks (sin PostgreSQL):

* JWT válido → permitido.
* Sin token / token inválido / expirado → 401.
* Autenticado sin permiso → 403 (cuando se implementen roles).
* Payload inválido → 400.
* Recurso inexistente → 404.

Comandos:

```text
go mod tidy
go test ./...
```

## OpenCode

Antes de modificar:

1. Leer este archivo y el `AGENTS.md` de la raíz.
2. Leer `../docs/api-contract.md` y `../docs/authentication.md`.
3. Leer `../docs/decisions/ADR-002-inventory-service-go.md`.
4. Revisar el código existente y reutilizar la estructura.
5. Hacer el cambio mínimo y agregar tests.
6. Ejecutar `go test ./...`.

No reescribir el módulo completo salvo requerimiento explícito.
