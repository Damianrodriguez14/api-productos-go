# API REST de Productos con Go y PostgreSQL

API REST completa para gestión de productos, conectada a PostgreSQL en Supabase.

## Tecnologías
- Go
- PostgreSQL (Supabase)
- API REST con endpoints CRUD

## Endpoints
- GET /productos → obtener todos los productos
- GET /productos/{id} → obtener un producto
- POST /productos → crear producto
- PUT /productos/{id} → actualizar producto
- DELETE /productos/{id} → eliminar producto

## Configuración
1. Clonar el repositorio
2. Crear archivo `.env` con:
DATABASE_URL=tu_url_de_supabase
3. Ejecutar con `go run main.go`

## Ejemplo de uso
POST /productos
{
  "nombre": "Remera básica",
  "precio": 15000,
  "stock": 50
}