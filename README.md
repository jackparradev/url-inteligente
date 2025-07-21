# URL Inteligente – Acortador de URLs en Go

Este proyecto implementa un **acortador de URLs** desarrollado en Go, cumpliendo con restricciones específicas como el uso exclusivo de almacenamiento en memoria, librerías estándar, y una arquitectura modular.

---

## 🚀 Funcionalidad

- Acorta URLs largas en un formato corto único.
- Redirecciona automáticamente desde la URL corta hacia la original.
- Persistencia en memoria sin uso de bases de datos reales.
- Manejo seguro de concurrencia.
- Validación básica de URLs y manejo de errores HTTP.

---

## 📁 Estructura del Proyecto

url-inteligente/
├── go.mod
├── main.go               # Punto de entrada
├── shorten.json          # Archivo de prueba opcional
├── internal/
│   ├── config/           # Configuración (pendiente de uso)
│   ├── handler/          # Endpoints HTTP
│   ├── service/          # Lógica de negocio (shortener y storage)
│   └── util/             # Funciones auxiliares

**Motivo de la estructura:**  
Usamos un enfoque tipo "clean architecture" adaptado a Go, separando claramente la lógica de negocio (`service`), la interacción HTTP (`handler`) y utilidades comunes (`util`). Esto mejora la mantenibilidad, pruebas unitarias y escalabilidad futura.

---

## 🔐 Generación de Códigos Cortos y Manejo de Colisiones

- El algoritmo usa un **hash SHA1** de la URL original junto con la marca de tiempo y el intento actual.
- Se toman los primeros 7 caracteres del hash para formar el código corto.
- En caso de colisión (el código ya existe), se reintenta hasta 5 veces agregando aleatoriedad adicional.
- No se utilizan librerías externas, solo `crypto/sha1`, `time` y `math/rand` de la librería estándar.

---

## 🔁 Redirección: HTTP 301 vs 307

Se utiliza **HTTP 301 (Moved Permanently)** porque:

- La relación entre una URL corta y la URL original es **permanente** una vez creada.
- Permite que navegadores y motores de búsqueda **cachen** la redirección, mejorando el rendimiento.
- 307 se evita, ya que está diseñado para redirecciones temporales o para mantener el método HTTP (como en POST → POST), lo cual no aplica aquí.

---

## 🧠 Concurrencia y Almacenamiento

- El almacenamiento se implementa mediante un `map[string]string` protegido con `sync.RWMutex`.
- Se permite acceso concurrente seguro para lecturas múltiples (`RLock`) y bloqueos exclusivos para escritura (`Lock`).
- No hay uso de bases de datos externas ni archivos: el almacenamiento es **puramente en memoria**.

---

## 🧪 Endpoints Principales

- `POST /shorten`: Acorta una URL (espera JSON `{ "url": "https://ejemplo.com" }`)
- `GET /{codigo}`: Redirige hacia la URL original asociada al código.

---

## ✅ Requisitos Cumplidos

- ✔ Uso exclusivo de `net/http` (sin frameworks complejos).
- ✔ Generación de códigos corta robusta y sin librerías externas.
- ✔ Almacenamiento concurrente en memoria con `sync.RWMutex`.
- ✔ Manejo exhaustivo de errores con códigos HTTP apropiados.
- ✔ Modularidad clara y estructura mantenible.

---

## 📦 Cómo ejecutar

```bash
go run main.go
