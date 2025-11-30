# SpecialDates

[![Stack: Go + Svelte + MySQL](https://img.shields.io/badge/Stack-Go%20%7C%20Svelte%20%7C%20MySQL-blue)](https://github.com/)
[![License: MIT](https://img.shields.io/badge/License-MIT-lightgrey)]

Aplicación web para la gestión de **fechas especiales** (cumpleaños, aniversarios, reuniones, etc.).  
Proyecto académico que demuestra la implementación de servicios web **REST** y **SOAP**, con arquitectura por capas (Presentación → Handlers → Service → Repository → DB).

---

## 📌 Resumen rápido

- **Backend:** Go  
- **Frontend:** Svelte + Vite  
- **Base de datos:** MySQL  
- **Funcionalidad:** CRUD (Create, Read, Update, Delete) de eventos/fechas, con consumo mediante REST y un endpoint SOAP didáctico.  
- **URL local por defecto:**  
  - Backend: `http://localhost:{puerto}`  
  - Frontend: `http://localhost:{puerto}`

---

## Características

- Interfaz en Svelte con formulario y tarjetas de eventos.  
- API REST completa (`/api/dates`) para CRUD.  
- Endpoint SOAP (`/soap`) para demostrar interoperabilidad con XML.  
- Arquitectura en capas: handlers, service, repository, db, models.  
- Validaciones (ej.: título obligatorio, no fechas en el pasado, etc.).

---

## Requisitos

- Go 1.25+  
- Node.js (LTS recomendado) y npm  
- MySQL (ej.: XAMPP, MySQL Server, MySQL Workbench)  
- Git (opcional)

---

## 📁 Estructura del repositorio
specialdates/
├─ backend/
│ ├─ go.mod
│ ├─ main.go
│ └─ internal/
│ ├─ db/db.go
│ ├─ models/date.go
│ ├─ repository/repository.go
│ ├─ service/service.go
│ └─ handlers/handlers.go
└─ frontend/
├─ package.json
├─ vite.config.js
├─ index.html
└─ src/
├─ main.js
└─ App.svelte


---

## 🔧 Crear la base de datos (MySQL)

Conéctate como administrador (`root`) y ejecuta:

```sql
CREATE DATABASE IF NOT EXISTS specialdates
  CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE specialdates;

CREATE TABLE IF NOT EXISTS dates (
  id INT AUTO_INCREMENT PRIMARY KEY,
  title VARCHAR(255) NOT NULL,
  description TEXT,
  event_datetime DATETIME NOT NULL,
  date_type VARCHAR(100),
  recurring TINYINT(1) NOT NULL DEFAULT 0,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

Crear usuario recomendado (localhost):

CREATE USER IF NOT EXISTS 'edit-username'@'localhost' IDENTIFIED BY 'edit-password';
GRANT ALL PRIVILEGES ON specialdates.* TO 'edit-username'@'localhost';
FLUSH PRIVILEGES;

```
⚙️ Variables de entorno

La app backend usa la variable DB_DSN con la cadena de conexión MySQL.

PowerShell (Windows): $env:DB_DSN = "edit-username:edit-password@tcp(localhost:{puerto})/specialdates?parseTime=true"
Linux / macOS (bash/zsh): export DB_DSN="edit-username:edit-password@tcp(localhost:{puerto})/specialdates?parseTime=true"

▶️ Ejecutar el backend (Go)

Desde la carpeta backend/: 
# instalar dependencias (si se requiere)
go mod tidy

# ejecutar
go run ./...

Salida esperada:
Servidor backend escuchando en :{puerto}

▶️ Ejecutar el frontend (Svelte + Vite)

Desde la carpeta frontend/:

npm install
npm run dev

Abre http://localhost:{puerto} en el navegador.

Si hay conflictos de dependencias npm (peer deps), borra node_modules y package-lock.json y vuelve a npm install. Se recomienda usar las versiones compatibles indicadas en package.json del repo.

📡 Endpoints principales (REST)

GET /api/dates — Listar todas las fechas.

GET /api/dates/{id} — Obtener fecha por id.

POST /api/dates — Crear fecha (JSON).

PUT /api/dates/{id} — Actualizar.

DELETE /api/dates/{id} — Eliminar.

Ejemplo JSON (POST):

{
  "title": "John's Birthday",
  "description": "Family dinner",
  "event_at": "2025-12-01T10:00:00Z",
  "date_type": "Birthday",
  "recurring": true
}

Ejemplo con curl:
curl -X POST http://localhost:{puerto}/api/dates \
  -H "Content-Type: application/json" \
  -d '{"title":"Demo","description":"Creada via curl","event_at":"2025-12-01T10:00:00Z","date_type":"Holiday","recurring":false}'

PowerShell (Invoke-RestMethod):
Invoke-RestMethod -Uri "http://localhost:{puerto}/api/dates" -Method Post -Headers @{ "Content-Type" = "application/json" } -Body '{"title":"Demo","description":"Creada via PowerShell","event_at":"2025-12-01T10:00:00Z","date_type":"Holiday","recurring":false}'

📄 Endpoint SOAP 

POST /soap — Recibe XML SOAP con operación CreateTask y responde CreateTaskResponse con <id>.

Ejemplo de body SOAP:

<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <CreateTask>
      <title>SOAP test</title>
      <description>Creada via SOAP</description>
    </CreateTask>
  </soap:Body>
</soap:Envelope>

✅ Verificar en la base de datos

En MySQL Workbench o consola:
USE specialdates;
SELECT id, title, event_datetime, date_type, recurring, created_at
FROM dates ORDER BY id DESC LIMIT 20;

🛠 Troubleshooting

DB_DSN no definida: exporta DB_DSN en la misma terminal antes de go run.

Error MySQL 1044 (Access denied): concede permisos al usuario o usa root para crear la BD y permisos.

Errores npm / versión de Svelte/Vite: ajustar package.json a versiones compatibles (el repo trae versiones recomendadas).

