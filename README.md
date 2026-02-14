# notes-api

![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)

A simple REST API to manage notes (id, title, content).

## Description

Notes API provides CRUD operations for notes over HTTP/JSON. It runs an in-memory store and exposes endpoints to create, read, update, and delete notes.

**Endpoints:**
- `GET /notes` — List all notes
- `GET /notes/:id` — Get a note by ID
- `POST /notes` — Create a note
- `PUT /notes/:id` — Update a note
- `DELETE /notes/:id` — Delete a note

## Run

```bash
go run .
```

Server runs on `http://localhost:8080`.

---

## TODO

- [ ] **A. Validation** — Reject empty title or too-long fields; return 400 with a clear message.
- [ ] **B. Split handlers** — Move handlers into a separate package (e.g. `handlers`) or files (e.g. `handlers.go`, `store.go`) and keep `main.go` small.
- [ ] **C. Logging** — Use `log` to log each request (method, path, status code).
- [ ] **D. File persistence** — On startup, load notes from a JSON file; on each create/update/delete, write the store back to the file.
